//go:build acceptance || image || images

package v2

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/clients"
	"github.com/vnpaycloud-console/gophercloud/v2/internal/acceptance/tools"
	"github.com/vnpaycloud-console/gophercloud/v2/openstack/image/v2/images"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
)

func TestImagesListEachPage(t *testing.T) {
	client, err := clients.NewImageV2Client()
	th.AssertNoErr(t, err)

	listOpts := images.ListOpts{
		Limit: 1,
	}

	pager := images.List(client, listOpts)
	err = pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		images, err := images.ExtractImages(page)
		if err != nil {
			t.Fatalf("Unable to extract images: %v", err)
		}

		for _, image := range images {
			tools.PrintResource(t, image)
			tools.PrintResource(t, image.Properties)
		}

		return true, nil
	})
	th.AssertNoErr(t, err)
}

func TestImagesListAllPages(t *testing.T) {
	client, err := clients.NewImageV2Client()
	th.AssertNoErr(t, err)

	image, err := CreateEmptyImage(t, client)
	th.AssertNoErr(t, err)
	defer DeleteImage(t, client, image)

	listOpts := images.ListOpts{}

	allPages, err := images.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allImages, err := images.ExtractImages(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, i := range allImages {
		tools.PrintResource(t, i)
		tools.PrintResource(t, i.Properties)

		if i.Name == image.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestImagesListByDate(t *testing.T) {
	client, err := clients.NewImageV2Client()
	th.AssertNoErr(t, err)

	date := time.Date(2014, 1, 1, 1, 1, 1, 0, time.UTC)
	listOpts := images.ListOpts{
		Limit: 1,
		CreatedAtQuery: &images.ImageDateQuery{
			Date:   date,
			Filter: images.FilterGTE,
		},
	}

	allPages, err := images.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allImages, err := images.ExtractImages(allPages)
	th.AssertNoErr(t, err)

	if len(allImages) == 0 {
		t.Fatalf("Query resulted in no results")
	}

	for _, image := range allImages {
		tools.PrintResource(t, image)
		tools.PrintResource(t, image.Properties)
	}

	date = time.Date(2049, 1, 1, 1, 1, 1, 0, time.UTC)
	listOpts = images.ListOpts{
		Limit: 1,
		CreatedAtQuery: &images.ImageDateQuery{
			Date:   date,
			Filter: images.FilterGTE,
		},
	}

	allPages, err = images.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allImages, err = images.ExtractImages(allPages)
	th.AssertNoErr(t, err)

	if len(allImages) > 0 {
		t.Fatalf("Expected 0 images, got %d", len(allImages))
	}
}

func TestImagesFilter(t *testing.T) {
	client, err := clients.NewImageV2Client()
	th.AssertNoErr(t, err)

	image, err := CreateEmptyImage(t, client)
	th.AssertNoErr(t, err)
	defer DeleteImage(t, client, image)

	listOpts := images.ListOpts{
		Tags:            []string{"foo", "bar"},
		ContainerFormat: "bare",
		DiskFormat:      "qcow2",
	}

	allPages, err := images.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allImages, err := images.ExtractImages(allPages)
	th.AssertNoErr(t, err)

	if len(allImages) == 0 {
		t.Fatalf("Query resulted in no results")
	}
}

func TestImagesUpdate(t *testing.T) {
	client, err := clients.NewImageV2Client()
	th.AssertNoErr(t, err)

	image, err := CreateEmptyImage(t, client)
	th.AssertNoErr(t, err)
	defer DeleteImage(t, client, image)

	newTags := []string{"foo", "bar"}

	updateOpts := images.UpdateOpts{
		images.ReplaceImageName{NewName: image.Name + "foo"},
		images.ReplaceImageTags{NewTags: newTags},
		images.ReplaceImageMinDisk{NewMinDisk: 21},
		images.ReplaceImageProtected{NewProtected: true},
		images.UpdateImageProperty{
			Op:    images.AddOp,
			Name:  "hw_disk_bus",
			Value: "scsi",
		},
		images.UpdateImageProperty{
			Op:   images.RemoveOp,
			Name: "architecture",
		},
	}

	newImage, err := images.Update(context.TODO(), client, image.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newImage)
	tools.PrintResource(t, newImage.Properties)

	th.AssertEquals(t, newImage.Name, image.Name+"foo")
	th.AssertEquals(t, newImage.Protected, true)

	sort.Strings(newTags)
	sort.Strings(newImage.Tags)
	th.AssertDeepEquals(t, newImage.Tags, newTags)

	// Because OpenStack is now adding additional properties automatically,
	// it's not possible to do an easy AssertDeepEquals.
	th.AssertEquals(t, newImage.Properties["hw_disk_bus"], "scsi")

	if _, ok := newImage.Properties["architecture"]; ok {
		t.Fatal("architecture property still exists")
	}

	// Now change image protection back to false or delete will fail
	updateOpts = images.UpdateOpts{
		images.ReplaceImageProtected{NewProtected: false},
	}
	_, err = images.Update(context.TODO(), client, image.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)
}
