package cli

import (
	"context"
	"testing"

	portainer "github.com/portainer/portainer/api"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
)

func Test_GetServiceAccount(t *testing.T) {

	t.Run("returns error if non-existent", func(t *testing.T) {
		k := &KubeClient{
			cli:        kfake.NewSimpleClientset(),
			instanceID: "test",
		}
		tokenData := &portainer.TokenData{ID: 1}
		_, err := k.GetPortainerUserServiceAccount(tokenData)
		if err == nil {
			t.Error("GetPortainerUserServiceAccount should fail with service account not found")
		}
	})

	t.Run("succeeds for cluster admin role", func(t *testing.T) {
		k := &KubeClient{
			cli:        kfake.NewSimpleClientset(),
			instanceID: "test",
		}

		tokenData := &portainer.TokenData{
			ID:       1,
			Role:     portainer.AdministratorRole,
			Username: portainerClusterAdminServiceAccountName,
		}
		serviceAccount := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: tokenData.Username,
			},
		}
		_, err := k.cli.CoreV1().ServiceAccounts(portainerNamespace).Create(context.Background(), serviceAccount, metav1.CreateOptions{})
		if err != nil {
			t.Errorf("failed to create service acount; err=%s", err)
		}
		defer k.cli.CoreV1().ServiceAccounts(portainerNamespace).Delete(context.Background(), serviceAccount.Name, metav1.DeleteOptions{})

		sa, err := k.GetPortainerUserServiceAccount(tokenData)
		if err != nil {
			t.Errorf("GetPortainerUserServiceAccount should succeed; err=%s", err)
		}

		want := "portainer-sa-clusteradmin"
		if sa.Name != want {
			t.Errorf("GetServiceAccount should succeed and return correct sa name; got=%s want=%s", sa.Name, want)
		}
	})

	t.Run("succeeds for standard user role", func(t *testing.T) {
		k := &KubeClient{
			cli:        kfake.NewSimpleClientset(),
			instanceID: "test",
		}

		tokenData := &portainer.TokenData{
			ID:   1,
			Role: portainer.StandardUserRole,
		}
		serviceAccountName := UserServiceAccountName(int(tokenData.ID), k.instanceID)
		serviceAccount := &v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name: serviceAccountName,
			},
		}
		_, err := k.cli.CoreV1().ServiceAccounts(portainerNamespace).Create(context.Background(), serviceAccount, metav1.CreateOptions{})
		if err != nil {
			t.Errorf("failed to create service acount; err=%s", err)
		}
		defer k.cli.CoreV1().ServiceAccounts(portainerNamespace).Delete(context.Background(), serviceAccount.Name, metav1.DeleteOptions{})

		sa, err := k.GetPortainerUserServiceAccount(tokenData)
		if err != nil {
			t.Errorf("GetPortainerUserServiceAccount should succeed; err=%s", err)
		}

		want := "portainer-sa-user-test-1"
		if sa.Name != want {
			t.Errorf("GetPortainerUserServiceAccount should succeed and return correct sa name; got=%s want=%s", sa.Name, want)
		}
	})

}
