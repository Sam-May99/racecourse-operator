/*
Copyright 2023.

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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	racecoursev1 "Sam-May99/racecourse-operator/api/v1"
)

// RacecourseApplicationReconciler reconciles a RacecourseApplication object
type RacecourseApplicationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=racecourse.my.domain,resources=racecourseapplications,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=racecourse.my.domain,resources=racecourseapplications/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=racecourse.my.domain,resources=racecourseapplications/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the RacecourseApplication object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *RacecourseApplicationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// Attempt to get the CR
	var racecourseApp racecoursev1.RacecourseApplication
	if err := r.Get(ctx, req.NamespacedName, &racecourseApp); err != nil {
		log.Log.Error(err, "Could not find the parent CR")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	for i := 0; i < int(*racecourseApp.Spec.Replicas); i++ {
		// Create the service fronting the traffic
		var racecourseAppService corev1.Service
		serviceName := fmt.Sprintf("%s-%s-%d", racecourseApp.Name, "svc", i)
		if err := r.Get(ctx, types.NamespacedName{Name: serviceName, Namespace: req.NamespacedName.Namespace}, &racecourseAppService); err != nil {

			service := r.createRacecourseService(&racecourseApp, serviceName)

			err = r.Create(context.TODO(), service)
			if err != nil {
				log.Log.Error(err, "Could not create service for the Racecourse")
				return ctrl.Result{}, err
			}
		}

		// Create the deployment for the pod
		var racecourseAppDeployment appsv1.Deployment
		deploymentName := fmt.Sprintf("%s-%s-%d", racecourseApp.Name, "deployment", i)
		if err := r.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: req.NamespacedName.Namespace}, &racecourseAppDeployment); err != nil {
			deployment := r.createRacecourseDeployment(&racecourseApp, deploymentName)
			racecourseApp.Status.Racecourses = append(racecourseApp.Status.Racecourses, deploymentName)

			err = r.Create(context.TODO(), deployment)
			if err != nil {
				log.Log.Error(err, "Could not create deployment for the Racecourse")
				return ctrl.Result{}, err
			}
		}
	}

	if err := r.Status().Update(ctx, &racecourseApp); err != nil {
		log.Log.Error(err, "Could not update the status of the racecourse")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// Creates the actual deployment for the Race Course
func (r *RacecourseApplicationReconciler) createRacecourseDeployment(racecourseApp *racecoursev1.RacecourseApplication, deploymentName string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": racecourseApp.Name,
			},
			Annotations: make(map[string]string),
			Name:        deploymentName,
			Namespace:   racecourseApp.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": racecourseApp.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": racecourseApp.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: racecourseApp.Spec.Containers,
				},
			},
		},
	}

	ctrl.SetControllerReference(racecourseApp, deployment, r.Scheme)
	return deployment
}

// Create the service resource for the Racecourse
func (r *RacecourseApplicationReconciler) createRacecourseService(racecourseApp *racecoursev1.RacecourseApplication, serviceName string) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": racecourseApp.Name,
			},
			Annotations: make(map[string]string),
			Name:        serviceName,
			Namespace:   racecourseApp.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": racecourseApp.Name,
			},
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       3000,
				TargetPort: intstr.FromInt(3000),
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	ctrl.SetControllerReference(racecourseApp, service, r.Scheme)
	return service
}

// SetupWithManager sets up the controller with the Manager.
func (r *RacecourseApplicationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&racecoursev1.RacecourseApplication{}).
		Complete(r)
}
