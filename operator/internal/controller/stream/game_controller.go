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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	streamv1 "indiegamestream.com/indiegamestream/api/stream/v1"
	stunnerv1 "indiegamestream.com/indiegamestream/api/stunner/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
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

	log.Info("Reconciling Game", "Name", game.Spec.Name, "ExecutableURL", game.Spec.ExecutableURL)

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

	//TODO: use stunner namespace
	var udpRoutes stunnerv1.UDPRouteList
	if err := r.List(ctx, &udpRoutes, client.InNamespace(req.Namespace), client.MatchingFields{udpRouteOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list child Jobs")
		return ctrl.Result{}, err
	}

	//show all udpRoutes
	/*
		for _, udpRoute := range udpRoutes.Items {
			log.Info("udpRoute", "Name", udpRoute.Name, "parentRefs", udpRoute.Spec.ParentRefs, "rules", udpRoute.Spec.Rules)
		}
	*/

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
		if false { //!reflect.DeepEqual(udpRoute.Spec, udpRouteRef.Spec) { //TODO: Add real check based on game spec
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

	// Check if the deployment already exists, if not create a new one
	found := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{Namespace: game.Namespace, Name: "coordinator-deployment"}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep, err := r.constructControllerDeploymentForGame(game)
		if err != nil {
			log.Error(err, "unable to construct deployment")
			return ctrl.Result{}, err
		}

		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "unable to create Deployment for Game", "game", game)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "unable to get Deployment for Game", "game", game)
		return ctrl.Result{}, err
	} else {
		//update deployment
		dep, err := r.constructControllerDeploymentForGame(game)
		if err != nil {
			log.Error(err, "unable to construct deployment")
			return ctrl.Result{}, err
		}

		//log.Info("Deployments", "Cur", found.Spec, "New", dep.Spec)

		// Update the deployment spec if necessary
		if false { //!reflect.DeepEqual(found.Spec, dep.Spec) { //TODO: Add real check based on game spec
			found.Spec = dep.Spec
			if err := r.Update(ctx, found); err != nil {
				log.Error(err, "unable to update Deployment for Game", "game", game)
				return ctrl.Result{}, err
			}
			log.V(1).Info("updated Deployment for Game", "game", game)
		}

	}

	// Check if the service already exists, if not create a new one
	foundSvc := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{Namespace: game.Namespace, Name: "coordinator-lb-svc"}, foundSvc)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		svc, err := r.constructControllerLoadBalancerForGame(game)
		if err != nil {
			log.Error(err, "unable to construct service")
			return ctrl.Result{}, err
		}

		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "unable to create Service for Game", "game", game)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "unable to get Service for Game", "game", game)
		return ctrl.Result{}, err
	} else {
		//update service
		svc, err := r.constructControllerLoadBalancerForGame(game)
		if err != nil {
			log.Error(err, "unable to construct service")
			return ctrl.Result{}, err
		}
		//log.Info("Services", "Cur", foundSvc.Spec, "New", svc.Spec)
		// Update the service spec if necessary
		if false { //!reflect.DeepEqual(foundSvc.Spec, svc.Spec) { //TODO: Add real check based on game spec
			foundSvc.Spec = svc.Spec
			if err := r.Update(ctx, foundSvc); err != nil {
				log.Error(err, "unable to update Service for Game", "game", game)
				return ctrl.Result{}, err

			}
			log.V(1).Info("updated Service for Game", "game", game)
		}
	}

	// Check if the deployment already exists, if not create a new one
	foundWorker := &appsv1.Deployment{}
	err = r.Get(ctx, client.ObjectKey{Namespace: game.Namespace, Name: "worker-deployment"}, foundWorker)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		dep, err := r.constructWorkerDeploymentForGame(game)
		if err != nil {
			log.Error(err, "unable to construct deployment")
			return ctrl.Result{}, err
		}

		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(err, "unable to create Deployment for Game", "game", game)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "unable to get Deployment for Game", "game", game)
		return ctrl.Result{}, err
	} else {
		//update deployment
		dep, err := r.constructWorkerDeploymentForGame(game)
		if err != nil {
			log.Error(err, "unable to construct deployment")
			return ctrl.Result{}, err
		}

		//log.Info("Deployments", "Cur", found.Spec, "New", dep.Spec)

		// Update the deployment spec if necessary
		if false { //!reflect.DeepEqual(foundWorker.Spec, dep.Spec) { //TODO: Add real check based on game spec
			foundWorker.Spec = dep.Spec
			if err := r.Update(ctx, foundWorker); err != nil {
				log.Error(err, "unable to update Deployment for Game", "game", game)
				return ctrl.Result{}, err
			}
			log.V(1).Info("updated Deployment for Game", "game", game)
		}

	}

	// Check if the service already exists, if not create a new one
	foundSvcWorker := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{Namespace: game.Namespace, Name: "worker-lb-svc"}, foundSvcWorker)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		svc, err := r.constructWorkerHTTPLoadBalancerForGame(game)
		if err != nil {
			log.Error(err, "unable to construct service")
			return ctrl.Result{}, err
		}

		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "unable to create Service for Game", "game", game)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "unable to get Service for Game", "game", game)
		return ctrl.Result{}, err
	} else {
		//update service
		svc, err := r.constructWorkerHTTPLoadBalancerForGame(game)
		if err != nil {
			log.Error(err, "unable to construct service")
			return ctrl.Result{}, err
		}
		//log.Info("Services", "Cur", foundSvc.Spec, "New", svc.Spec)
		// Update the service spec if necessary
		if false { //!reflect.DeepEqual(foundSvcWorker.Spec, svc.Spec) { //TODO: Add real check based on game spec
			foundSvcWorker.Spec = svc.Spec
			if err := r.Update(ctx, foundSvcWorker); err != nil {
				log.Error(err, "unable to update Service for Game", "game", game)
				return ctrl.Result{}, err

			}
			log.V(1).Info("updated Service for Game", "game", game)
		}
	}

	// Check if the service already exists, if not create a new one
	foundSvcWorkerUDP := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{Namespace: game.Namespace, Name: "worker-ci-udp-svc"}, foundSvcWorkerUDP)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		svc, err := r.constructWorkerUDPLoadBalancerForGame(game)
		if err != nil {
			log.Error(err, "unable to construct service")
			return ctrl.Result{}, err
		}

		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			log.Error(err, "unable to create Service for Game", "game", game)
			return ctrl.Result{}, err
		}

	} else if err != nil {
		log.Error(err, "unable to get Service for Game", "game", game)
		return ctrl.Result{}, err
	} else {
		//update service
		svc, err := r.constructWorkerUDPLoadBalancerForGame(game)
		if err != nil {
			log.Error(err, "unable to construct service")
			return ctrl.Result{}, err
		}
		//log.Info("Services", "Cur", foundSvc.Spec, "New", svc.Spec)
		// Update the service spec if necessary
		if false { //!reflect.DeepEqual(foundSvcWorkerUDP.Spec, svc.Spec) { //TODO: Add real check based on game spec
			foundSvcWorkerUDP.Spec = svc.Spec
			if err := r.Update(ctx, foundSvcWorkerUDP); err != nil {
				log.Error(err, "unable to update Service for Game", "game", game)
				return ctrl.Result{}, err

			}
			log.V(1).Info("updated Service for Game", "game", game)
		}
	}

	// Finally, we update the status block of the Game resource to reflect the current state of the world
	// Note that Status is a subresource, so changes to it are ignored by the cache, hence the need to update it manually
	//game.Status.Nodes = nodes
	//game.Status.Phase = phase
	//TODO: Add nginx ingress url to status
	if err := r.Status().Update(ctx, game); err != nil {
		log.Error(err, "unable to update Game status")
		return ctrl.Result{}, err
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
func int32Ptr(i int32) *int32 {
	return &i
}

func (r *GameReconciler) constructControllerDeploymentForGame(game *streamv1.Game) (*appsv1.Deployment, error) {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "coordinator-deployment",
			Namespace: game.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "coordinator"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "coordinator"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "coordinator",
							Image:   "valniae/snekyrepo:crdi",
							Command: []string{"coordinator"},
							Args:    []string{"--v=5"},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8000,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "CLOUD_GAME_WEBRTC_ICESERVERS_0_CREDENTIAL",
									Value: "TODO_ADD_CREDS",
								},
								{
									Name:  "CLOUD_GAME_WEBRTC_ICESERVERS_0_URL",
									Value: "turn:10.0.0.1:3478", //TODO: should be loadbalancer IP
								},
								{
									Name:  "CLOUD_GAME_WEBRTC_ICESERVERS_0_USERNAME",
									Value: "gilroy",
								},
							}, //TODO mount game executable and config game config file
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(game, dep, r.Scheme); err != nil {
		return nil, err
	}

	return dep, nil
}

func (r *GameReconciler) constructWorkerDeploymentForGame(game *streamv1.Game) (*appsv1.Deployment, error) {
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "worker-deployment",
			Namespace: game.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "worker"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "worker"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "worker",
							Image:   "valniae/snekyrepo:crdi",
							Command: []string{"worker"},
							Args:    []string{"--v=5"},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9000,
								},
								{
									ContainerPort: 8443,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "CLOUD_GAME_EMULATOR_AUTOSAVESEC",
									Value: "3",
								},
							}, //TODO mount game executable and config game config file
						},
					},
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(game, dep, r.Scheme); err != nil {
		return nil, err
	}

	return dep, nil
}

func (r *GameReconciler) constructControllerLoadBalancerForGame(game *streamv1.Game) (*corev1.Service, error) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "coordinator-lb-svc",
			Namespace: game.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "coordinator"},
			Ports: []corev1.ServicePort{
				{
					Port:       8000,
					TargetPort: intstr.FromInt(8000),
				},
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}

	if err := ctrl.SetControllerReference(game, svc, r.Scheme); err != nil {
		return nil, err
	}

	return svc, nil
}

func (r *GameReconciler) constructWorkerHTTPLoadBalancerForGame(game *streamv1.Game) (*corev1.Service, error) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "worker-lb-svc",
			Namespace: game.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "worker"},
			Ports: []corev1.ServicePort{
				{
					Port:       9000,
					TargetPort: intstr.FromInt(9000),
				},
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}

	if err := ctrl.SetControllerReference(game, svc, r.Scheme); err != nil {
		return nil, err
	}

	return svc, nil
}

func (r *GameReconciler) constructWorkerUDPLoadBalancerForGame(game *streamv1.Game) (*corev1.Service, error) {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "worker-ci-udp-svc",
			Namespace: game.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": "worker"},
			Ports: []corev1.ServicePort{
				{
					Port:       8443,
					TargetPort: intstr.FromInt(8443),
				},
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}

	if err := ctrl.SetControllerReference(game, svc, r.Scheme); err != nil {
		return nil, err
	}

	return svc, nil
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
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
