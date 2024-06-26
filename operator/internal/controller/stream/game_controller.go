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
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
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
//+kubebuilder:rbac:groups=stunner.l7mp.io,resources=gatewayconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;

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

	log.Info("Reconciling Game", "Name", game.Spec.Name, "FileName", game.Spec.FileName)

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

			if err := r.deleteExternalResources(ctx, game); err != nil {
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
	workerName := fmt.Sprintf("worker-lb-svc-%s", game.Name)
	coordinatorName := fmt.Sprintf("coordinator-lb-svc-%s", game.Name)
	udpRouteName := fmt.Sprintf("udproute-%s", game.Name)
	deploymentCoordName := fmt.Sprintf("deployment-coord-%s", game.Name)
	deploymentWorkerName := fmt.Sprintf("deployment-worker-%s", game.Name)
	workerUDPName := fmt.Sprintf("worker-ci-udp-svc-%s", game.Name)

	result, err := r.ensureResource(ctx, game, "UDPRoute", udpRouteName, game.Namespace, workerUDPName)
	if err != nil {
		return result, err
	}

	// Get data for controller deployment creation

	gatewayConfig := &stunnerv1.GatewayConfig{}
	err = r.Get(ctx, client.ObjectKey{
		Namespace: "stunner",
		Name:      "stunner-gatewayconfig",
	}, gatewayConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to get GatewayConfig: %v", err))
	}

	fmt.Printf("Username: %s\n", gatewayConfig.Spec.UserName)
	fmt.Printf("Password: %s\n", gatewayConfig.Spec.Password)

	gatewayIP, err := waitForLoadBalancerIP(ctx, r.Client, "stunner", "udp-gateway")
	if err != nil {
		log.Error(err, "unable to get LoadBalancer IP for Gateway")
		return ctrl.Result{}, err
	}
	log.Info("Gateway LoadBalancer IP", "IP", gatewayIP)

	result, err = r.ensureResource(ctx, game, "Deployment-Coordinator", deploymentCoordName, game.Namespace, gatewayConfig, gatewayIP)
	if err != nil {
		return result, err
	}

	result, err = r.ensureResource(ctx, game, "Service", coordinatorName, game.Namespace, "coordinator", int32(8000))
	if err != nil {
		return result, err
	}

	result, err = r.ensureResource(ctx, game, "Service", workerName, game.Namespace, "worker", int32(9000))
	if err != nil {
		return result, err
	}

	coordIP, err := waitForLoadBalancerIP(ctx, r.Client, game.Namespace, coordinatorName)
	if err != nil {
		log.Error(err, "unable to get LoadBalancer IP for Coordinator")
		return ctrl.Result{}, err
	}
	log.Info("Coordinator LoadBalancer IP", "IP", coordIP)

	workerIP, err := waitForLoadBalancerIP(ctx, r.Client, game.Namespace, workerName)

	if err != nil {
		log.Error(err, "unable to get LoadBalancer IP for Worker")
		return ctrl.Result{}, err
	}
	log.Info("Worker LoadBalancer IP", "IP", workerIP)

	result, err = r.ensureResource(ctx, game, "Deployment-Worker", deploymentWorkerName, game.Namespace, coordIP, workerIP)
	if err != nil {
		return result, err
	}

	result, err = r.ensureResource(ctx, game, "Service-UDP", workerUDPName, game.Namespace, "worker", int32(8443))
	if err != nil {
		return result, err
	}

	// Finally, we update the status block of the Game resource to reflect the current state of the world
	// Note that Status is a subresource, so changes to it are ignored by the cache, hence the need to update it manually
	//game.Status.Nodes = nodes
	//game.Status.Phase = phase
	outsidehostname, err := waitForLoadBalancerHostname(ctx, r.Client, game.Namespace, coordinatorName)
	if err != nil {
		log.Error(err, "unable to get LoadBalancer Hostname for Coordinator")
		return ctrl.Result{}, err
	}

	game.Status.URL = fmt.Sprintf("http://%s", outsidehostname)
	//TODO: Add nginx ingress url to status
	if err := r.Status().Update(ctx, game); err != nil {
		log.Error(err, "unable to update Game status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *GameReconciler) ensureResource(ctx context.Context, game *streamv1.Game, resourceType string, resourceName string, resourceNamespace string, args ...interface{}) (ctrl.Result, error) {
	var log = log.FromContext(ctx)
	var resource client.Object
	var constructFunc func(*streamv1.Game, string, ...interface{}) (client.Object, error)

	switch resourceType {
	case "UDPRoute":
		resource = &stunnerv1.UDPRoute{}
		resourceNamespace = "stunner"
		constructFunc = func(g *streamv1.Game, resourceName string, params ...interface{}) (client.Object, error) {
			if len(params) != 1 {
				return nil, fmt.Errorf("invalid number of arguments for Deployment: expected 1, got %d", len(params))
			}
			serviceRef, ok := params[0].(string)
			if !ok {
				return nil, fmt.Errorf("invalid parameter type for GatewayConfig")
			}
			udproute, err := r.constructUDPRouteForGame(g, resourceName, serviceRef)
			return udproute, err
		}
	case "Deployment-Coordinator":
		resource = &appsv1.Deployment{}
		constructFunc = func(g *streamv1.Game, resourceName string, params ...interface{}) (client.Object, error) {
			if len(params) != 2 {
				return nil, fmt.Errorf("invalid number of arguments for Deployment: expected 2, got %d", len(params))
			}
			gatewayConfig, ok := params[0].(*stunnerv1.GatewayConfig)
			if !ok {
				return nil, fmt.Errorf("invalid parameter type for GatewayConfig")
			}
			gatewayIP, ok := params[1].(string)
			if !ok {
				return nil, fmt.Errorf("invalid parameter type for gatewayIP")
			}
			return r.constructControllerDeploymentForGame(g, resourceName, gatewayConfig, gatewayIP)
		}
	case "Deployment-Worker":
		resource = &appsv1.Deployment{}
		constructFunc = func(g *streamv1.Game, resourceName string, params ...interface{}) (client.Object, error) {
			if len(params) != 2 {
				return nil, fmt.Errorf("invalid number of arguments for Deployment: expected 2, got %d", len(params))
			}
			coordIP, ok := params[0].(string)
			if !ok {
				return nil, fmt.Errorf("invalid parameter type for coordIP")
			}
			workerIP, ok := params[1].(string)
			if !ok {
				return nil, fmt.Errorf("invalid parameter type for gatewayIP")
			}
			return r.constructWorkerDeploymentForGame(g, resourceName, coordIP, workerIP)
		}
	case "Service":
		resource = &corev1.Service{}
		constructFunc = func(g *streamv1.Game, resourceName string, params ...interface{}) (client.Object, error) {
			if len(params) != 2 {
				return nil, fmt.Errorf("invalid number of arguments for service: expected 2, got %d", len(params))
			}
			label, ok := params[0].(string)
			if !ok {
				return nil, fmt.Errorf("invalid parameter type for label")
			}
			port, ok := params[1].(int32)
			if !ok {
				return nil, fmt.Errorf("invalid parameter type for port")
			}
			udproute, err := r.constructLoadBalancer(g, resourceName, label, port)

			return udproute, err
		}
	case "Service-UDP":
		resource = &corev1.Service{}
		constructFunc = func(g *streamv1.Game, resourceName string, params ...interface{}) (client.Object, error) {
			if len(params) != 2 {
				return nil, fmt.Errorf("invalid number of arguments for service: expected 2, got %d", len(params))
			}
			label, ok := params[0].(string)
			if !ok {
				return nil, fmt.Errorf("invalid parameter type for label")
			}
			port, ok := params[1].(int32)
			if !ok {
				return nil, fmt.Errorf("invalid parameter type for port")
			}
			udproute, err := r.constructLoadBalancerUDP(g, resourceName, label, port)

			return udproute, err
		}

	// Add more cases as needed for different resource types
	default:
		return ctrl.Result{}, fmt.Errorf("unsupported resource type %s", resourceType)
	}

	err := r.Get(ctx, client.ObjectKey{Namespace: resourceNamespace, Name: resourceName}, resource)
	if err != nil && errors.IsNotFound(err) {
		newResource, err := constructFunc(game, resourceName, args...)
		if err != nil {
			log.Error(err, "unable to construct resource", "type", resourceType)
			return ctrl.Result{}, err
		}

		log.Info("Creating a new resource", "type", resourceType, "namespace", newResource.GetNamespace(), "name", newResource.GetName())
		if err = r.Create(ctx, newResource); err != nil {
			log.Error(err, "unable to create resource for Game", "game", game)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "unable to get resource for Game", "game", game)
		return ctrl.Result{}, err
	} else {
		// TODO: handle updates
	}

	return ctrl.Result{}, nil
}

func (r *GameReconciler) deleteExternalResources(ctx context.Context, game *streamv1.Game) error {
	//manually delete udproute

	udpRoute := &stunnerv1.UDPRoute{}
	err := r.Get(ctx, client.ObjectKey{Namespace: "stunner", Name: fmt.Sprintf("udproute-%s", game.Name)}, udpRoute)
	if err != nil {
		return nil
	}
	if err := r.Delete(ctx, udpRoute); err != nil {
		return err
	}

	return nil
}

func (r *GameReconciler) constructUDPRouteForGame(game *streamv1.Game, resourceName string, serviceRef string) (*stunnerv1.UDPRoute, error) {
	// We want job names for a given nominal start time to have a deterministic name to avoid the same job being created twice

	udpRoute := &stunnerv1.UDPRoute{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
			Name:        resourceName,
			Namespace:   "stunner",
		},
		Spec: stunnerv1.UDPRouteSpec{
			ParentRefs: []stunnerv1.ParentRefSpec{
				{
					Name:      "udp-gateway",
					Namespace: "stunner",
				},
			},
			Rules: []stunnerv1.RulesSpec{
				{
					BackendRefs: []stunnerv1.BackendRefSpec{
						{
							Name:      serviceRef,
							Namespace: game.Namespace,
						},
					},
				},
			},
		},
	}

	//if err := ctrl.SetControllerReference(game, udpRoute, r.Scheme); err != nil {
	//	return nil, err
	//}

	return udpRoute, nil
}
func int32Ptr(i int32) *int32 {
	return &i
}

func (r *GameReconciler) constructControllerDeploymentForGame(game *streamv1.Game, resourceName string, gatewayConfig *stunnerv1.GatewayConfig, gatewayIP string) (*appsv1.Deployment, error) {
	fullpath := fmt.Sprintf("/usr/local/share/cloud-game/assets/games/%s", game.Spec.FileName)
	newSelector := fmt.Sprintf("%s-%s", "coordinator", game.Name)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: game.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": newSelector},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": newSelector},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "coordinator",
							Image:   "ghcr.io/giongto35/cloud-game/cloud-game:v3.0.5",
							Command: []string{"./coordinator"},
							Args:    []string{""},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8000,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "CLOUD_GAME_WEBRTC_ICESERVERS_0_CREDENTIAL",
									Value: gatewayConfig.Spec.Password,
								},
								{
									Name:  "CLOUD_GAME_WEBRTC_ICESERVERS_0_URLS",
									Value: fmt.Sprintf("turn:%s:3478", gatewayIP),
								},
								{
									Name:  "CLOUD_GAME_WEBRTC_ICESERVERS_0_USERNAME",
									Value: gatewayConfig.Spec.UserName,
								},
								{
									Name:  "CLOUD_GAME_WEBRTC_ICESERVERS_1_CREDENTIAL",
									Value: gatewayConfig.Spec.Password,
								},
								{
									Name:  "CLOUD_GAME_WEBRTC_ICESERVERS_1_URLS",
									Value: fmt.Sprintf("turn:%s:3478", gatewayIP),
								},
								{
									Name:  "CLOUD_GAME_WEBRTC_ICESERVERS_1_USERNAME",
									Value: gatewayConfig.Spec.UserName,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "gamestorage",
									MountPath: fullpath,
									SubPath:   game.Name,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "gamestorage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "azure-blob-pvc",
								},
							},
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

func (r *GameReconciler) constructWorkerDeploymentForGame(game *streamv1.Game, resourceName string, coordIP string, workerIP string) (*appsv1.Deployment, error) {
	fullpath := fmt.Sprintf("/usr/local/share/cloud-game/assets/games/%s", game.Spec.FileName)
	newSelector := fmt.Sprintf("%s-%s", "worker", game.Name)
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: game.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": newSelector},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": newSelector},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "worker",
							Image:   "ghcr.io/giongto35/cloud-game/cloud-game:v3.0.5",
							Command: []string{"./worker"},
							Args:    []string{""},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8443,
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "CLOUD_GAME_EMULATOR_AUTOSAVESEC",
									Value: "3",
								},
								{
									Name:  "CLOUD_GAME_WORKER_NETWORK_COORDINATORADDRESS",
									Value: fmt.Sprintf("%s:80", coordIP),
								},
								{
									Name:  "CLOUD_GAME_WORKER_NETWORK_PUBLICADDRESS",
									Value: workerIP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "gamestorage",
									MountPath: fullpath,
									SubPath:   game.Name,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "gamestorage",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: "azure-blob-pvc",
								},
							},
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

func (r *GameReconciler) constructLoadBalancer(game *streamv1.Game, name string, selector string, port int32) (*corev1.Service, error) {

	className := "tailscale"
	annotation := fmt.Sprintf("%s-%s", game.Spec.Name, game.Name)
	newSelector := fmt.Sprintf("%s-%s", selector, game.Name)
	annotations := map[string]string{}
	outsidePort := port
	if strings.HasPrefix(selector, "coordinator") {
		annotations["tailscale.com/hostname"] = annotation
		outsidePort = 80
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   game.Namespace,
			Annotations: annotations,
		},
		Spec: corev1.ServiceSpec{
			Selector:          map[string]string{"app": newSelector},
			LoadBalancerClass: &className,
			Ports: []corev1.ServicePort{
				{
					Port:       outsidePort,
					TargetPort: intstr.FromInt32(port),
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
func (r *GameReconciler) constructLoadBalancerUDP(game *streamv1.Game, name string, selector string, port int32) (*corev1.Service, error) {
	newSelector := fmt.Sprintf("%s-%s", "worker", game.Name)
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: game.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": newSelector},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolUDP,
					Port:       port,
					TargetPort: intstr.FromInt32(port),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	if err := ctrl.SetControllerReference(game, svc, r.Scheme); err != nil {
		return nil, err
	}

	return svc, nil
}

func waitForLoadBalancerIP(ctx context.Context, k8sClient client.Client, namespace, serviceName string) (string, error) {
	var ip string

	err := wait.PollUntilContextTimeout(ctx, 5*time.Second, 20*time.Second, true, func(ctx context.Context) (bool, error) {
		svc := &corev1.Service{}
		if err := k8sClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: serviceName}, svc); err != nil {
			return false, err
		}
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			ip = svc.Status.LoadBalancer.Ingress[0].IP
			if ip != "" {
				return true, nil
			}
			if len(svc.Status.LoadBalancer.Ingress) > 1 {
				ip = svc.Status.LoadBalancer.Ingress[1].IP
				if ip != "" {
					return true, nil
				}
			}
		}
		return false, nil
	})

	if err != nil {
		return "", err
	}
	return ip, nil
}

func waitForLoadBalancerHostname(ctx context.Context, k8sClient client.Client, namespace, serviceName string) (string, error) {
	var hostname string

	err := wait.PollUntilContextTimeout(ctx, 5*time.Second, 20*time.Second, true, func(ctx context.Context) (bool, error) {
		svc := &corev1.Service{}
		if err := k8sClient.Get(ctx, client.ObjectKey{Namespace: namespace, Name: serviceName}, svc); err != nil {
			return false, err
		}
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			hostname = svc.Status.LoadBalancer.Ingress[0].Hostname
			if hostname != "" {
				return true, nil
			}
			if len(svc.Status.LoadBalancer.Ingress) > 1 {
				hostname = svc.Status.LoadBalancer.Ingress[1].Hostname
				if hostname != "" {
					return true, nil
				}
			}
		}
		return false, nil
	})

	//fallback on ip
	if err != nil {
		ip, err := waitForLoadBalancerIP(ctx, k8sClient, namespace, serviceName)
		if err != nil {
			return "", err
		}
		return ip, nil
	}

	return hostname, nil
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
		//Owns(&stunnerv1.UDPRoute{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Complete(r)
}
