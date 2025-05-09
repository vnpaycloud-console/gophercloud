package nodes

import (
	"context"
	"fmt"

	"github.com/vnpaycloud-console/gophercloud/v2"
	"github.com/vnpaycloud-console/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToNodeListQuery() (string, error)
	ToNodeListDetailQuery() (string, error)
}

// Provision state reports the current provision state of the node, these are only used in filtering
type ProvisionState string

const (
	Enroll       ProvisionState = "enroll"
	Verifying    ProvisionState = "verifying"
	Manageable   ProvisionState = "manageable"
	Available    ProvisionState = "available"
	Active       ProvisionState = "active"
	DeployWait   ProvisionState = "wait call-back"
	Deploying    ProvisionState = "deploying"
	DeployFail   ProvisionState = "deploy failed"
	DeployDone   ProvisionState = "deploy complete"
	DeployHold   ProvisionState = "deploy hold"
	Deleting     ProvisionState = "deleting"
	Deleted      ProvisionState = "deleted"
	Cleaning     ProvisionState = "cleaning"
	CleanWait    ProvisionState = "clean wait"
	CleanFail    ProvisionState = "clean failed"
	CleanHold    ProvisionState = "clean hold"
	Error        ProvisionState = "error"
	Rebuild      ProvisionState = "rebuild"
	Inspecting   ProvisionState = "inspecting"
	InspectFail  ProvisionState = "inspect failed"
	InspectWait  ProvisionState = "inspect wait"
	Adopting     ProvisionState = "adopting"
	AdoptFail    ProvisionState = "adopt failed"
	Rescue       ProvisionState = "rescue"
	RescueFail   ProvisionState = "rescue failed"
	Rescuing     ProvisionState = "rescuing"
	UnrescueFail ProvisionState = "unrescue failed"
	RescueWait   ProvisionState = "rescue wait"
	Unrescuing   ProvisionState = "unrescuing"
	Servicing    ProvisionState = "servicing"
	ServiceWait  ProvisionState = "service wait"
	ServiceFail  ProvisionState = "service failed"
	ServiceHold  ProvisionState = "service hold"
)

// TargetProvisionState is used when setting the provision state for a node.
type TargetProvisionState string

const (
	TargetActive   TargetProvisionState = "active"
	TargetDeleted  TargetProvisionState = "deleted"
	TargetManage   TargetProvisionState = "manage"
	TargetProvide  TargetProvisionState = "provide"
	TargetInspect  TargetProvisionState = "inspect"
	TargetAbort    TargetProvisionState = "abort"
	TargetClean    TargetProvisionState = "clean"
	TargetAdopt    TargetProvisionState = "adopt"
	TargetRescue   TargetProvisionState = "rescue"
	TargetUnrescue TargetProvisionState = "unrescue"
	TargetRebuild  TargetProvisionState = "rebuild"
	TargetService  TargetProvisionState = "service"
	TargetUnhold   TargetProvisionState = "unhold"
)

const (
	StepHold     string = "hold"
	StepWait     string = "wait"
	StepPowerOn  string = "power_on"
	StepPowerOff string = "power_off"
	StepReboot   string = "reboot"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the node attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	// Filter the list by specific instance UUID
	InstanceUUID string `q:"instance_uuid"`

	// Filter the list by chassis UUID
	ChassisUUID string `q:"chassis_uuid"`

	// Filter the list by maintenance set to True or False
	Maintenance bool `q:"maintenance"`

	// Nodes which are, or are not, associated with an instance_uuid.
	Associated bool `q:"associated"`

	// Only return those with the specified provision_state.
	ProvisionState ProvisionState `q:"provision_state"`

	// Filter the list with the specified driver.
	Driver string `q:"driver"`

	// Filter the list with the specified resource class.
	ResourceClass string `q:"resource_class"`

	// Filter the list with the specified conductor_group.
	ConductorGroup string `q:"conductor_group"`

	// Filter the list with the specified fault.
	Fault string `q:"fault"`

	// One or more fields to be returned in the response.
	Fields []string `q:"fields" format:"comma-separated"`

	// Requests a page size of items.
	Limit int `q:"limit"`

	// The ID of the last-seen item.
	Marker string `q:"marker"`

	// Sorts the response by the requested sort direction.
	SortDir string `q:"sort_dir"`

	// Sorts the response by the this attribute value.
	SortKey string `q:"sort_key"`

	// A string or UUID of the tenant who owns the baremetal node.
	Owner string `q:"owner"`
}

// ToNodeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToNodeListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list nodes accessible to you.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToNodeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return NodePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ToNodeListDetailQuery formats a ListOpts into a query string for the list details API.
func (opts ListOpts) ToNodeListDetailQuery() (string, error) {
	// Detail endpoint can't filter by Fields
	if len(opts.Fields) > 0 {
		return "", fmt.Errorf("fields is not a valid option when getting a detailed listing of nodes")
	}

	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Return a list of bare metal Nodes with complete details. Some filtering is possible by passing in flags in ListOpts,
// but you cannot limit by the fields returned.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	// This URL is deprecated. In the future, we should compare the microversion and if >= 1.43, hit the listURL
	// with ListOpts{Detail: true,}
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToNodeListDetailQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return NodePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get requests details on a single node, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToNodeCreateMap() (map[string]any, error)
}

// CreateOpts specifies node creation parameters.
type CreateOpts struct {
	// The interface to configure automated cleaning for a Node.
	// Requires microversion 1.47 or later.
	AutomatedClean *bool `json:"automated_clean,omitempty"`

	// The BIOS interface for a Node, e.g. “redfish”.
	BIOSInterface string `json:"bios_interface,omitempty"`

	// The boot interface for a Node, e.g. “pxe”.
	BootInterface string `json:"boot_interface,omitempty"`

	// The conductor group for a node. Case-insensitive string up to 255 characters, containing a-z, 0-9, _, -, and ..
	ConductorGroup string `json:"conductor_group,omitempty"`

	// The console interface for a node, e.g. “no-console”.
	ConsoleInterface string `json:"console_interface,omitempty"`

	// The deploy interface for a node, e.g. “iscsi”.
	DeployInterface string `json:"deploy_interface,omitempty"`

	// All the metadata required by the driver to manage this Node. List of fields varies between drivers, and can
	// be retrieved from the /v1/drivers/<DRIVER_NAME>/properties resource.
	DriverInfo map[string]any `json:"driver_info,omitempty"`

	// name of the driver used to manage this Node.
	Driver string `json:"driver,omitempty"`

	// A set of one or more arbitrary metadata key and value pairs.
	Extra map[string]any `json:"extra,omitempty"`

	// The firmware interface for a node, e.g. "redfish"
	FirmwareInterface string `json:"firmware_interface,omitempty"`

	// The interface used for node inspection, e.g. “no-inspect”.
	InspectInterface string `json:"inspect_interface,omitempty"`

	// Interface for out-of-band node management, e.g. “ipmitool”.
	ManagementInterface string `json:"management_interface,omitempty"`

	// Human-readable identifier for the Node resource. May be undefined. Certain words are reserved.
	Name string `json:"name,omitempty"`

	// Which Network Interface provider to use when plumbing the network connections for this Node.
	NetworkInterface string `json:"network_interface,omitempty"`

	// Interface used for performing power actions on the node, e.g. “ipmitool”.
	PowerInterface string `json:"power_interface,omitempty"`

	// Physical characteristics of this Node. Populated during inspection, if performed. Can be edited via the REST
	// API at any time.
	Properties map[string]any `json:"properties,omitempty"`

	// Interface used for configuring RAID on this node, e.g. “no-raid”.
	RAIDInterface string `json:"raid_interface,omitempty"`

	// The interface used for node rescue, e.g. “no-rescue”.
	RescueInterface string `json:"rescue_interface,omitempty"`

	// A string which can be used by external schedulers to identify this Node as a unit of a specific type
	// of resource.
	ResourceClass string `json:"resource_class,omitempty"`

	// Interface used for attaching and detaching volumes on this node, e.g. “cinder”.
	StorageInterface string `json:"storage_interface,omitempty"`

	// The UUID for the resource.
	UUID string `json:"uuid,omitempty"`

	// Interface for vendor-specific functionality on this node, e.g. “no-vendor”.
	VendorInterface string `json:"vendor_interface,omitempty"`

	// A string or UUID of the tenant who owns the baremetal node.
	Owner string `json:"owner,omitempty"`

	// Static network configuration to use during deployment and cleaning.
	NetworkData map[string]any `json:"network_data,omitempty"`

	// Whether disable_power_off is enabled or disabled on this node.
	// Requires microversion 1.95 or later.
	DisablePowerOff *bool `json:"disable_power_off,omitempty"`
}

// ToNodeCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToNodeCreateMap() (map[string]any, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Create requests a node to be created
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToNodeCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, createURL(client), reqBody, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type Patch interface {
	ToNodeUpdateMap() (map[string]any, error)
}

// UpdateOpts is a slice of Patches used to update a node
type UpdateOpts []Patch

type UpdateOp string

const (
	ReplaceOp UpdateOp = "replace"
	AddOp     UpdateOp = "add"
	RemoveOp  UpdateOp = "remove"
)

type UpdateOperation struct {
	Op    UpdateOp `json:"op" required:"true"`
	Path  string   `json:"path" required:"true"`
	Value any      `json:"value,omitempty"`
}

func (opts UpdateOperation) ToNodeUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Update requests that a node be updated
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	body := make([]map[string]any, len(opts))
	for i, patch := range opts {
		result, err := patch.ToNodeUpdateMap()
		if err != nil {
			r.Err = err
			return
		}

		body[i] = result
	}
	resp, err := client.Patch(ctx, updateURL(client, id), body, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete requests that a node be removed
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Request that Ironic validate whether the Node’s driver has enough information to manage the Node. This polls each
// interface on the driver, and returns the status of that interface.
func Validate(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ValidateResult) {
	resp, err := client.Get(ctx, validateURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Inject NMI (Non-Masking Interrupts) for the given Node. This feature can be used for hardware diagnostics, and
// actual support depends on a driver.
func InjectNMI(ctx context.Context, client *gophercloud.ServiceClient, id string) (r InjectNMIResult) {
	resp, err := client.Put(ctx, injectNMIURL(client, id), map[string]string{}, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type BootDeviceOpts struct {
	BootDevice string `json:"boot_device"` // e.g., 'pxe', 'disk', etc.
	Persistent bool   `json:"persistent"`  // Whether this is one-time or not
}

// BootDeviceOptsBuilder allows extensions to add additional parameters to the
// SetBootDevice request.
type BootDeviceOptsBuilder interface {
	ToBootDeviceMap() (map[string]any, error)
}

// ToBootDeviceSetMap assembles a request body based on the contents of a BootDeviceOpts.
func (opts BootDeviceOpts) ToBootDeviceMap() (map[string]any, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Set the boot device for the given Node, and set it persistently or for one-time boot. The exact behaviour
// of this depends on the hardware driver.
func SetBootDevice(ctx context.Context, client *gophercloud.ServiceClient, id string, bootDevice BootDeviceOptsBuilder) (r SetBootDeviceResult) {
	reqBody, err := bootDevice.ToBootDeviceMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, bootDeviceURL(client, id), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get the current boot device for the given Node.
func GetBootDevice(ctx context.Context, client *gophercloud.ServiceClient, id string) (r BootDeviceResult) {
	resp, err := client.Get(ctx, bootDeviceURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Retrieve the acceptable set of supported boot devices for a specific Node.
func GetSupportedBootDevices(ctx context.Context, client *gophercloud.ServiceClient, id string) (r SupportedBootDeviceResult) {
	resp, err := client.Get(ctx, supportedBootDeviceURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// An interface type for a deploy (or clean) step.
type StepInterface string

const (
	InterfaceBIOS       StepInterface = "bios"
	InterfaceDeploy     StepInterface = "deploy"
	InterfaceFirmware   StepInterface = "firmware"
	InterfaceManagement StepInterface = "management"
	InterfacePower      StepInterface = "power"
	InterfaceRAID       StepInterface = "raid"
)

// A cleaning step has required keys ‘interface’ and ‘step’, and optional key ‘args’. If specified,
// the value for ‘args’ is a keyword variable argument dictionary that is passed to the cleaning step
// method.
type CleanStep struct {
	Interface StepInterface  `json:"interface" required:"true"`
	Step      string         `json:"step" required:"true"`
	Args      map[string]any `json:"args,omitempty"`
}

// A service step looks the same as a cleaning step.
type ServiceStep = CleanStep

// A deploy step has required keys ‘interface’, ‘step’, ’args’ and ’priority’.
// The value for ‘args’ is a keyword variable argument dictionary that is passed to the deploy step
// method. Priority is a numeric priority at which the step is running.
type DeployStep struct {
	Interface StepInterface  `json:"interface" required:"true"`
	Step      string         `json:"step" required:"true"`
	Args      map[string]any `json:"args" required:"true"`
	Priority  int            `json:"priority" required:"true"`
}

// ProvisionStateOptsBuilder allows extensions to add additional parameters to the
// ChangeProvisionState request.
type ProvisionStateOptsBuilder interface {
	ToProvisionStateMap() (map[string]any, error)
}

// Starting with Ironic API version 1.56, a configdrive may be a JSON object with structured data.
// Prior to this version, it must be a base64-encoded, gzipped ISO9660 image.
type ConfigDrive struct {
	MetaData    map[string]any `json:"meta_data,omitempty"`
	NetworkData map[string]any `json:"network_data,omitempty"`
	UserData    any            `json:"user_data,omitempty"`
}

// ProvisionStateOpts for a request to change a node's provision state. A config drive should be base64-encoded
// gzipped ISO9660 image. Deploy steps are supported starting with API 1.69.
type ProvisionStateOpts struct {
	Target         TargetProvisionState `json:"target" required:"true"`
	ConfigDrive    any                  `json:"configdrive,omitempty"`
	CleanSteps     []CleanStep          `json:"clean_steps,omitempty"`
	DeploySteps    []DeployStep         `json:"deploy_steps,omitempty"`
	ServiceSteps   []ServiceStep        `json:"service_steps,omitempty"`
	RescuePassword string               `json:"rescue_password,omitempty"`
}

// ToProvisionStateMap assembles a request body based on the contents of a CreateOpts.
func (opts ProvisionStateOpts) ToProvisionStateMap() (map[string]any, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Request a change to the Node’s provision state. Acceptable target states depend on the Node’s current provision
// state. More detailed documentation of the Ironic State Machine is available in the developer docs.
func ChangeProvisionState(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ProvisionStateOptsBuilder) (r ChangeStateResult) {
	reqBody, err := opts.ToProvisionStateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, provisionStateURL(client, id), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type TargetPowerState string

// TargetPowerState is used when changing the power state of a node.
const (
	PowerOn       TargetPowerState = "power on"
	PowerOff      TargetPowerState = "power off"
	Rebooting     TargetPowerState = "rebooting"
	SoftPowerOff  TargetPowerState = "soft power off"
	SoftRebooting TargetPowerState = "soft rebooting"
)

// PowerStateOptsBuilder allows extensions to add additional parameters to the ChangePowerState request.
type PowerStateOptsBuilder interface {
	ToPowerStateMap() (map[string]any, error)
}

// PowerStateOpts for a request to change a node's power state.
type PowerStateOpts struct {
	Target  TargetPowerState `json:"target" required:"true"`
	Timeout int              `json:"timeout,omitempty"`
}

// ToPowerStateMap assembles a request body based on the contents of a PowerStateOpts.
func (opts PowerStateOpts) ToPowerStateMap() (map[string]any, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Request to change a Node's power state.
func ChangePowerState(ctx context.Context, client *gophercloud.ServiceClient, id string, opts PowerStateOptsBuilder) (r ChangePowerStateResult) {
	reqBody, err := opts.ToPowerStateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, powerStateURL(client, id), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// This is the desired RAID configuration on the bare metal node.
type RAIDConfigOpts struct {
	LogicalDisks []LogicalDisk `json:"logical_disks"`
}

// RAIDConfigOptsBuilder allows extensions to modify a set RAID config request.
type RAIDConfigOptsBuilder interface {
	ToRAIDConfigMap() (map[string]any, error)
}

// RAIDLevel type is used to specify the RAID level for a logical disk.
type RAIDLevel string

const (
	RAID0  RAIDLevel = "0"
	RAID1  RAIDLevel = "1"
	RAID2  RAIDLevel = "2"
	RAID5  RAIDLevel = "5"
	RAID6  RAIDLevel = "6"
	RAID10 RAIDLevel = "1+0"
	RAID50 RAIDLevel = "5+0"
	RAID60 RAIDLevel = "6+0"
	JBOD   RAIDLevel = "JBOD"
)

// DiskType is used to specify the disk type for a logical disk, e.g. hdd or ssd.
type DiskType string

const (
	HDD DiskType = "hdd"
	SSD DiskType = "ssd"
)

// InterfaceType is used to specify the interface for a logical disk.
type InterfaceType string

const (
	SATA InterfaceType = "sata"
	SCSI InterfaceType = "scsi"
	SAS  InterfaceType = "sas"
)

type LogicalDisk struct {
	// Size (Integer) of the logical disk to be created in GiB.  If unspecified, "MAX" will be used.
	SizeGB *int `json:"size_gb"`

	// RAID level for the logical disk.
	RAIDLevel RAIDLevel `json:"raid_level" required:"true"`

	// Name of the volume. Should be unique within the Node. If not specified, volume name will be auto-generated.
	VolumeName string `json:"volume_name,omitempty"`

	// Set to true if this is the root volume. At most one logical disk can have this set to true.
	IsRootVolume *bool `json:"is_root_volume,omitempty"`

	// Set to true if this logical disk can share physical disks with other logical disks.
	SharePhysicalDisks *bool `json:"share_physical_disks,omitempty"`

	// If this is not specified, disk type will not be a criterion to find backing physical disks
	DiskType DiskType `json:"disk_type,omitempty"`

	// If this is not specified, interface type will not be a criterion to find backing physical disks.
	InterfaceType InterfaceType `json:"interface_type,omitempty"`

	// Integer, number of disks to use for the logical disk. Defaults to minimum number of disks required
	// for the particular RAID level.
	NumberOfPhysicalDisks int `json:"number_of_physical_disks,omitempty"`

	// The name of the controller as read by the RAID interface.
	Controller string `json:"controller,omitempty"`

	// A list of physical disks to use as read by the RAID interface.
	PhysicalDisks []any `json:"physical_disks,omitempty"`
}

func (opts RAIDConfigOpts) ToRAIDConfigMap() (map[string]any, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if body["logical_disks"] != nil {
		for _, v := range body["logical_disks"].([]any) {
			if logicalDisk, ok := v.(map[string]any); ok {
				if logicalDisk["size_gb"] == nil {
					logicalDisk["size_gb"] = "MAX"
				}
			}
		}
	}

	return body, nil
}

// Request to change a Node's RAID config.
func SetRAIDConfig(ctx context.Context, client *gophercloud.ServiceClient, id string, raidConfigOptsBuilder RAIDConfigOptsBuilder) (r ChangeStateResult) {
	reqBody, err := raidConfigOptsBuilder.ToRAIDConfigMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, raidConfigURL(client, id), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListBIOSSettingsOptsBuilder allows extensions to add additional parameters to the
// ListBIOSSettings request.
type ListBIOSSettingsOptsBuilder interface {
	ToListBIOSSettingsOptsQuery() (string, error)
}

// ListBIOSSettingsOpts defines query options that can be passed to ListBIOSettings
type ListBIOSSettingsOpts struct {
	// Provide additional information for the BIOS Settings
	Detail bool `q:"detail"`

	// One or more fields to be returned in the response.
	Fields []string `q:"fields" format:"comma-separated"`
}

// ToListBIOSSettingsOptsQuery formats a ListBIOSSettingsOpts into a query string
func (opts ListBIOSSettingsOpts) ToListBIOSSettingsOptsQuery() (string, error) {
	if opts.Detail && len(opts.Fields) > 0 {
		return "", fmt.Errorf("cannot have both fields and detail options for BIOS settings")
	}

	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get the current BIOS Settings for the given Node.
// To use the opts requires microversion 1.74.
func ListBIOSSettings(ctx context.Context, client *gophercloud.ServiceClient, id string, opts ListBIOSSettingsOptsBuilder) (r ListBIOSSettingsResult) {
	url := biosListSettingsURL(client, id)
	if opts != nil {

		query, err := opts.ToListBIOSSettingsOptsQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	resp, err := client.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get one BIOS Setting for the given Node.
func GetBIOSSetting(ctx context.Context, client *gophercloud.ServiceClient, id string, setting string) (r GetBIOSSettingResult) {
	resp, err := client.Get(ctx, biosGetSettingURL(client, id, setting), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CallVendorPassthruOpts defines query options that can be passed to any VendorPassthruCall
type CallVendorPassthruOpts struct {
	Method string `q:"method"`
}

// ToGetSubscriptionMap assembles a query based on the contents of a CallVendorPassthruOpts
func ToGetAllSubscriptionMap(opts CallVendorPassthruOpts) (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Get all vendor_passthru methods available for the given Node.
func GetVendorPassthruMethods(ctx context.Context, client *gophercloud.ServiceClient, id string) (r VendorPassthruMethodsResult) {
	resp, err := client.Get(ctx, vendorPassthruMethodsURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get all subscriptions available for the given Node.
func GetAllSubscriptions(ctx context.Context, client *gophercloud.ServiceClient, id string, method CallVendorPassthruOpts) (r GetAllSubscriptionsVendorPassthruResult) {
	query, err := ToGetAllSubscriptionMap(method)
	if err != nil {
		r.Err = err
		return
	}
	url := vendorPassthruCallURL(client, id) + query
	resp, err := client.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// The desired subscription id on the baremetal node.
type GetSubscriptionOpts struct {
	Id string `json:"id"`
}

// ToGetSubscriptionMap assembles a query based on the contents of CallVendorPassthruOpts and a request body based on the contents of a GetSubscriptionOpts
func ToGetSubscriptionMap(method CallVendorPassthruOpts, opts GetSubscriptionOpts) (string, map[string]any, error) {
	q, err := gophercloud.BuildQueryString(method)
	if err != nil {
		return q.String(), nil, err
	}
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return q.String(), nil, err
	}

	return q.String(), body, nil
}

// Get a subscription on the given Node.
func GetSubscription(ctx context.Context, client *gophercloud.ServiceClient, id string, method CallVendorPassthruOpts, subscriptionOpts GetSubscriptionOpts) (r SubscriptionVendorPassthruResult) {
	query, reqBody, err := ToGetSubscriptionMap(method, subscriptionOpts)
	if err != nil {
		r.Err = err
		return
	}
	url := vendorPassthruCallURL(client, id) + query
	resp, err := client.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{
		JSONBody: reqBody,
		OkCodes:  []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// The desired subscription to be deleted from the baremetal node.
type DeleteSubscriptionOpts struct {
	Id string `json:"id"`
}

// ToDeleteSubscriptionMap assembles a query based on the contents of CallVendorPassthruOpts and a request body based on the contents of a DeleteSubscriptionOpts
func ToDeleteSubscriptionMap(method CallVendorPassthruOpts, opts DeleteSubscriptionOpts) (string, map[string]any, error) {
	q, err := gophercloud.BuildQueryString(method)
	if err != nil {
		return q.String(), nil, err
	}
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return q.String(), nil, err
	}
	return q.String(), body, nil
}

// Delete a subscription on the given node.
func DeleteSubscription(ctx context.Context, client *gophercloud.ServiceClient, id string, method CallVendorPassthruOpts, subscriptionOpts DeleteSubscriptionOpts) (r DeleteSubscriptionVendorPassthruResult) {
	query, reqBody, err := ToDeleteSubscriptionMap(method, subscriptionOpts)
	if err != nil {
		r.Err = err
		return
	}
	url := vendorPassthruCallURL(client, id) + query
	resp, err := client.Delete(ctx, url, &gophercloud.RequestOpts{
		JSONBody: reqBody,
		OkCodes:  []int{200, 202, 204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return r
}

// The desired subscription to be created from the baremetal node.
type CreateSubscriptionOpts struct {
	Destination string              `json:"Destination"`
	EventTypes  []string            `json:"EventTypes,omitempty"`
	HttpHeaders []map[string]string `json:"HttpHeaders,omitempty"`
	Context     string              `json:"Context,omitempty"`
	Protocol    string              `json:"Protocol,omitempty"`
}

// ToCreateSubscriptionMap assembles a query based on the contents of CallVendorPassthruOpts and a request body based on the contents of a CreateSubscriptionOpts
func ToCreateSubscriptionMap(method CallVendorPassthruOpts, opts CreateSubscriptionOpts) (string, map[string]any, error) {
	q, err := gophercloud.BuildQueryString(method)
	if err != nil {
		return q.String(), nil, err
	}
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return q.String(), nil, err
	}
	return q.String(), body, nil
}

// Creates a subscription on the given node.
func CreateSubscription(ctx context.Context, client *gophercloud.ServiceClient, id string, method CallVendorPassthruOpts, subscriptionOpts CreateSubscriptionOpts) (r SubscriptionVendorPassthruResult) {
	query, reqBody, err := ToCreateSubscriptionMap(method, subscriptionOpts)
	if err != nil {
		r.Err = err
		return
	}
	url := vendorPassthruCallURL(client, id) + query
	resp, err := client.Post(ctx, url, reqBody, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return r
}

// MaintenanceOpts for a request to set the node's maintenance mode.
type MaintenanceOpts struct {
	Reason string `json:"reason,omitempty"`
}

// MaintenanceOptsBuilder allows extensions to add additional parameters to the SetMaintenance request.
type MaintenanceOptsBuilder interface {
	ToMaintenanceMap() (map[string]any, error)
}

// ToMaintenanceMap assembles a request body based on the contents of a MaintenanceOpts.
func (opts MaintenanceOpts) ToMaintenanceMap() (map[string]any, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Request to set the Node's maintenance mode.
func SetMaintenance(ctx context.Context, client *gophercloud.ServiceClient, id string, opts MaintenanceOptsBuilder) (r SetMaintenanceResult) {
	reqBody, err := opts.ToMaintenanceMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, maintenanceURL(client, id), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Request to unset the Node's maintenance mode.
func UnsetMaintenance(ctx context.Context, client *gophercloud.ServiceClient, id string) (r SetMaintenanceResult) {
	resp, err := client.Delete(ctx, maintenanceURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetInventory return stored data from successful inspection.
func GetInventory(ctx context.Context, client *gophercloud.ServiceClient, id string) (r InventoryResult) {
	resp, err := client.Get(ctx, inventoryURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListFirmware return the list of Firmware components for the given Node.
func ListFirmware(ctx context.Context, client *gophercloud.ServiceClient, id string) (r ListFirmwareResult) {
	resp, err := client.Get(ctx, firmwareListURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type VirtualMediaDeviceType string

const (
	VirtualMediaDisk   VirtualMediaDeviceType = "disk"
	VirtualMediaCD     VirtualMediaDeviceType = "cdrom"
	VirtualMediaFloppy VirtualMediaDeviceType = "floppy"
)

type ImageDownloadSource string

const (
	ImageDownloadSourceHTTP  ImageDownloadSource = "http"
	ImageDownloadSourceLocal ImageDownloadSource = "local"
	ImageDownloadSourceSwift ImageDownloadSource = "swift"
)

// The desired virtual media attachment on the baremetal node.
type AttachVirtualMediaOpts struct {
	DeviceType          VirtualMediaDeviceType `json:"device_type"`
	ImageURL            string                 `json:"image_url"`
	ImageDownloadSource ImageDownloadSource    `json:"image_download_source,omitempty"`
}

type AttachVirtualMediaOptsBuilder interface {
	ToAttachVirtualMediaMap() (map[string]any, error)
}

func (opts AttachVirtualMediaOpts) ToAttachVirtualMediaMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Request to attach a virtual media device to the Node.
func AttachVirtualMedia(ctx context.Context, client *gophercloud.ServiceClient, id string, opts AttachVirtualMediaOptsBuilder) (r VirtualMediaAttachResult) {
	reqBody, err := opts.ToAttachVirtualMediaMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, virtualMediaURL(client, id), reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// The desired virtual media detachment on the baremetal node.
type DetachVirtualMediaOpts struct {
	DeviceTypes []VirtualMediaDeviceType `q:"device_types" format:"comma-separated"`
}

type DetachVirtualMediaOptsBuilder interface {
	ToDetachVirtualMediaOptsQuery() (string, error)
}

func (opts DetachVirtualMediaOpts) ToDetachVirtualMediaOptsQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Request to detach a virtual media device from the Node.
func DetachVirtualMedia(ctx context.Context, client *gophercloud.ServiceClient, id string, opts DetachVirtualMediaOptsBuilder) (r VirtualMediaDetachResult) {
	query, err := opts.ToDetachVirtualMediaOptsQuery()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Delete(ctx, virtualMediaURL(client, id)+query, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
