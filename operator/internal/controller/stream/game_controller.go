/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	streamv1 "indiegamestream.com/indiegamestream/api/stream/v1"
	stunnerv1 "indiegamestream.com/indiegamestream/api/stunner/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// GameReconciler reconciles a Game object
type GameReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=stream.indiegamestream.com,resources=games,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=stream.indiegamestream.com,resources=games/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=stream.indiegamestream.com,resources=games/finalizers,verbs=update
//+kubebuilder:rbac:groups=stunner.l7mp.io,resources=udproutes,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Game object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *GameReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var log = log.FromContext(ctx)

	log.Info("Request", "Incoming", req)

	game := &streamv1.Game{}
	if err := r.Get(ctx, req.NamespacedName, game); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Game resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch Game")
		return ctrl.Result{}, err
	}

	// name of our custom finalizer
	gameFinalizer := "game.stream.indiegamestream.com/finalizer"

	// examine DeletionTimestamp to determine if object is under deletion
	if game.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// to registering our finalizer.
		if !controllerutil.ContainsFinalizer(game, gameFinalizer) {
			controllerutil.AddFinalizer(game, gameFinalizer)
			log.Info("Finalizer added", "Name", game.Name, "Finalizer", gameFinalizer)
			if err := r.Update(ctx, game); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(game, gameFinalizer) {
			// our finalizer is present, so lets handle any external dependency
			log.Info("Game is being deleted", "Name", game.Name)

			if err := r.deleteExternalResources(game); err != nil {
				log.Error(err, "Error deleting external ressources")
				// if fail to delete the external dependency here, return with error
				// so that it can be retried.
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(game, gameFinalizer)
			if err := r.Update(ctx, game); err != nil {
				log.Error(err, "Error updating game after removing finalizer", "Name", game.Name)
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	log.Info("Reconciling Game", "Name", game.Spec.Name, "ExecutableURL", game.Spec.ExecutableURL)

	//TODO: use stunner namespace
	var udpRoutes stunnerv1.UDPRouteList
	if err := r.List(ctx, &udpRoutes, client.InNamespace(req.Namespace), client.MatchingFields{udpRouteOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list child Jobs")
		return ctrl.Result{}, err
	}

	//show all udpRoutes
	for _, udpRoute := range udpRoutes.Items {
		log.Info("udpRoute", "Name", udpRoute.Name, "parentRefs", udpRoute.Spec.ParentRefs, "rules", udpRoute.Spec.Rules)
	}

	if len(udpRoutes.Items) == 0 {
		// No udpRoutes found, define a new one

		udpRoute, err := r.constructUDPRouteForGame(game)
		if err != nil {
			log.Error(err, "unable to construct udproute")
			// don't bother requeuing until we get a change to the spec
			return ctrl.Result{}, err
		}

		// ...and create it on the cluster
		if err := r.Create(ctx, udpRoute); err != nil {
			log.Error(err, "unable to create UDPRoute for Gamne", "game", game)
			return ctrl.Result{}, err
		}

		log.V(1).Info("created UDPRoute for Game", "game", game)

	} else if len(udpRoutes.Items) == 1 {
		// We have exactly one udpRoute, make sure it's up to date
		udpRoute := udpRoutes.Items[0]

		//Create reference object from game

		udpRouteRef, err := r.constructUDPRouteForGame(game)
		if err != nil {
			log.Error(err, "unable to construct udproute")
			// don't bother requeuing until we get a change to the spec
			return ctrl.Result{}, err
		}

		// Update the udpRoute spec if necessary
		if !reflect.DeepEqual(udpRoute.Spec, udpRouteRef.Spec) {
			udpRoute.Spec = udpRouteRef.Spec
			if err := r.Update(ctx, &udpRoute); err != nil {
				log.Error(err, "unable to update UDPRoute for Game", "game", game)
				return ctrl.Result{}, err
			}
			log.V(1).Info("updated UDPRoute for Game", "game", game)
		}
	} else {
		return ctrl.Result{}, fmt.Errorf("found multiple UDPRoute for the same Game %s/%s", game.Namespace, game.Name)
	}

	return ctrl.Result{}, nil
}

func (r *GameReconciler) deleteExternalResources(game *streamv1.Game) error {
	return nil
}

func (r *GameReconciler) constructUDPRouteForGame(game *streamv1.Game) (*stunnerv1.UDPRoute, error) {
	// We want job names for a given nominal start time to have a deterministic name to avoid the same job being created twice
	name := fmt.Sprintf("udproute-%s", game.Name)

	udpRoute := &stunnerv1.UDPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        name,
			Namespace:   game.Namespace,
		},
		Spec: stunnerv1.UDPRouteSpec{
			ParentRefs: []stunnerv1.ParentRefSpec{
				{
					Name:      "game-gateway",
					Namespace: game.Namespace,
				},
			},
			Rules: []stunnerv1.RulesSpec{
				{
					BackendRefs: []stunnerv1.BackendRefSpec{
						{
							Name:      "game-service-backend",
							Namespace: game.Namespace,
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(game, udpRoute, r.Scheme); err != nil {
		return nil, err
	}

	return udpRoute, nil
}

var (
	udpRouteOwnerKey = ".metadata.controller"
	apiGVStr         = streamv1.GroupVersion.String()
)

// SetupWithManager sets up the controller with the Manager.
func (r *GameReconciler) SetupWithManager(mgr ctrl.Manager) error {

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &stunnerv1.UDPRoute{}, udpRouteOwnerKey, func(rawObj client.Object) []string {
		// grab the job object, extract the owner...
		job := rawObj.(*stunnerv1.UDPRoute)
		owner := metav1.GetControllerOf(job)
		if owner == nil {
			return nil
		}
		// ...make sure it's a CronJob...
		if owner.APIVersion != apiGVStr || owner.Kind != "Game" {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&streamv1.Game{}).
		Owns(&stunnerv1.UDPRoute{}).
		Complete(r)
}
