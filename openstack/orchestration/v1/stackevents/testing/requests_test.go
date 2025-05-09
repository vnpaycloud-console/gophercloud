package testing

import (
	"context"
	"testing"

	"github.com/vnpaycloud-console/gophercloud/v2/openstack/orchestration/v1/stackevents"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
	th "github.com/vnpaycloud-console/gophercloud/v2/testhelper"
	fake "github.com/vnpaycloud-console/gophercloud/v2/testhelper/client"
)

func TestFindEvents(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleFindSuccessfully(t, FindOutput)

	actual, err := stackevents.Find(context.TODO(), fake.ServiceClient(), "postman_stack").Extract()
	th.AssertNoErr(t, err)

	expected := FindExpected
	th.AssertDeepEquals(t, expected, actual)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t, ListOutput)

	count := 0
	err := stackevents.List(fake.ServiceClient(), "hello_world", "49181cd6-169a-4130-9455-31185bbfc5bf", nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := stackevents.ExtractEvents(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ListExpected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListResourceEvents(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListResourceEventsSuccessfully(t, ListResourceEventsOutput)

	count := 0
	err := stackevents.ListResourceEvents(fake.ServiceClient(), "hello_world", "49181cd6-169a-4130-9455-31185bbfc5bf", "my_resource", nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := stackevents.ExtractResourceEvents(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ListResourceEventsExpected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestGetEvent(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t, GetOutput)

	actual, err := stackevents.Get(context.TODO(), fake.ServiceClient(), "hello_world", "49181cd6-169a-4130-9455-31185bbfc5bf", "my_resource", "93940999-7d40-44ae-8de4-19624e7b8d18").Extract()
	th.AssertNoErr(t, err)

	expected := GetExpected
	th.AssertDeepEquals(t, expected, actual)
}
