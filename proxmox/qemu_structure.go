package proxmox

import (
	"strconv"
	"strings"

	pxapi "github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var resourceQemuSchema = map[string]*schema.Schema{
	"name": {
		Type:     schema.TypeString,
		Required: true,
	},
	"desc": {
		Type:     schema.TypeString,
		Optional: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.TrimSpace(old) == strings.TrimSpace(new)
		},
	},
	"target_node": {
		Type:     schema.TypeString,
		Required: true,
	},
	"bios": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "seabios",
	},
	"onboot": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	},
	"boot": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "cdn",
	},
	"bootdisk": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"agent": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	"iso": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"clone": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	},
	"full_clone": {
		Type:     schema.TypeBool,
		Optional: true,
		ForceNew: true,
		Default:  true,
	},
	"hastate": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"qemu_os": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "l26",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			if new == "l26" {
				return len(d.Get("clone").(string)) > 0 // the cloned source may have a different os, which we shoud leave alone
			}
			return strings.TrimSpace(old) == strings.TrimSpace(new)
		},
	},
	"memory": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  512,
	},
	"balloon": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	"cores": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  1,
	},
	"sockets": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  1,
	},
	"vcpus": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  0,
	},
	"cpu": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "host",
	},
	"numa": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"hotplug": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "network,disk,usb",
	},
	"scsihw": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "",
	},
	"vga": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "std",
				},
				"memory": {
					Type:     schema.TypeInt,
					Optional: true,
				},
			},
		},
	},
	"network": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": &schema.Schema{
					Type:     schema.TypeInt,
					Required: true,
				},
				"model": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"macaddr": &schema.Schema{
					// TODO: Find a way to set MAC address in .tf config.
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
				},
				"bridge": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "nat",
				},
				"tag": &schema.Schema{
					Type:        schema.TypeInt,
					Optional:    true,
					Description: "VLAN tag.",
					Default:     -1,
				},
				"firewall": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"rate": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
					Default:  -1,
				},
				"queues": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
					Default:  -1,
				},
				"link_down": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
			},
		},
	},
	"disk": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": &schema.Schema{
					Type:     schema.TypeInt,
					Required: true,
				},
				"type": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"storage": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"storage_type": &schema.Schema{
					Type:        schema.TypeString,
					Optional:    true,
					Default:     "dir",
					Description: "One of PVE types as described: https://pve.proxmox.com/wiki/Storage",
				},
				"size": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
				"format": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "raw",
				},
				"cache": &schema.Schema{
					Type:     schema.TypeString,
					Optional: true,
					Default:  "none",
				},
				"backup": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"iothread": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"replicate": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"mbps": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"mbps_rd": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"mbps_rd_max": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"mbps_wr": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
				"mbps_wr_max": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
					Default:  0,
				},
			},
		},
	},
	"serial": &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": &schema.Schema{
					Type:     schema.TypeInt,
					Required: true,
				},
				"type": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	},
	"os_type": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"os_network_config": {
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.TrimSpace(old) == strings.TrimSpace(new)
		},
	},
	"force_create": {
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	},
	"clone_wait": {
		Type:     schema.TypeInt,
		Optional: true,
		Default:  15,
	},
	"ci_wait": { // how long to wait before provision
		Type:     schema.TypeInt,
		Optional: true,
		Default:  30,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			if old == "" {
				return true // old empty ok
			}
			return strings.TrimSpace(old) == strings.TrimSpace(new)
		},
	},
	"ciuser": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"cipassword": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"cicustom": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"searchdomain": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"nameserver": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"sshkeys": {
		Type:     schema.TypeString,
		Optional: true,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return strings.TrimSpace(old) == strings.TrimSpace(new)
		},
	},
	"ipconfig0": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"ipconfig1": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"ipconfig2": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"pool": {
		Type:     schema.TypeString,
		Optional: true,
	},
}

func flattenVmQemu(vmr *pxapi.VmRef, config *pxapi.ConfigQemu, d *schema.ResourceData) {
	d.SetId(resourceId(vmr.Node(), "qemu", vmr.VmId()))
	d.Set("target_node", vmr.Node())
	d.Set("name", config.Name)
	d.Set("desc", config.Description)
	d.Set("pool", config.Pool)
	d.Set("bios", config.Bios)
	d.Set("onboot", config.Onboot)
	d.Set("boot", config.Boot)
	d.Set("bootdisk", config.BootDisk)
	d.Set("agent", config.Agent)
	d.Set("memory", config.Memory)
	d.Set("balloon", config.Balloon)
	d.Set("cores", config.QemuCores)
	d.Set("sockets", config.QemuSockets)
	d.Set("vcpus", config.QemuVcpus)
	d.Set("cpu", config.QemuCpu)
	d.Set("numa", config.QemuNuma)
	d.Set("hotplug", config.Hotplug)
	d.Set("scsihw", config.Scsihw)
	d.Set("hastate", vmr.HaState())
	d.Set("qemu_os", config.QemuOs)
	// Cloud-init.
	d.Set("ciuser", config.CIuser)
	d.Set("cipassword", config.CIpassword)
	d.Set("cicustom", config.CIcustom)
	d.Set("searchdomain", config.Searchdomain)
	d.Set("nameserver", config.Nameserver)
	d.Set("sshkeys", config.Sshkeys)
	d.Set("ipconfig0", config.Ipconfig0)
	d.Set("ipconfig1", config.Ipconfig1)
	d.Set("ipconfig2", config.Ipconfig2)
	// Disks.
	configDisksSet := d.Get("disk").(*schema.Set)
	activeDisksSet := UpdateDevicesSet(configDisksSet, config.QemuDisks)
	d.Set("disk", activeDisksSet)
	// Display.
	activeVgaSet := d.Get("vga").(*schema.Set)
	if len(activeVgaSet.List()) > 0 {
		d.Set("features", UpdateDeviceConfDefaults(config.QemuVga, activeVgaSet))
	}
	// Networks.
	configNetworksSet := d.Get("network").(*schema.Set)
	activeNetworksSet := UpdateDevicesSet(configNetworksSet, config.QemuNetworks)
	d.Set("network", activeNetworksSet)
	d.Set("pool", vmr.Pool())
	//Serials
	configSerialsSet := d.Get("serial").(*schema.Set)
	activeSerialSet := UpdateDevicesSet(configSerialsSet, config.QemuSerials)
	d.Set("serial", activeSerialSet)
}

// Converting from schema.TypeSet to map of id and conf for each device,
// which will be sent to Proxmox API.
func DevicesSetToMap(devicesSet *schema.Set) pxapi.QemuDevices {

	devicesMap := pxapi.QemuDevices{}

	for _, set := range devicesSet.List() {
		setMap, isMap := set.(map[string]interface{})
		if isMap {
			setID := setMap["id"].(int)
			devicesMap[setID] = setMap
		}
	}
	return devicesMap
}

// Update schema.TypeSet with new values comes from Proxmox API.
// TODO: Maybe it's better to create a new Set instead add to current one.
func UpdateDevicesSet(
	devicesSet *schema.Set,
	devicesMap pxapi.QemuDevices,
) *schema.Set {

	configDevicesMap := DevicesSetToMap(devicesSet)

	activeDevicesMap := updateDevicesDefaults(devicesMap, configDevicesMap)

	for _, setConf := range devicesSet.List() {
		devicesSet.Remove(setConf)
		setConfMap := setConf.(map[string]interface{})
		deviceID := setConfMap["id"].(int)
		// Value type should be one of types allowed by Terraform schema types.
		for key, value := range activeDevicesMap[deviceID] {
			// This nested switch is used for nested config like in `net[n]`,
			// where Proxmox uses `key=<0|1>` in string" at the same time
			// a boolean could be used in ".tf" files.
			switch setConfMap[key].(type) {
			case bool:
				switch value.(type) {
				// If the key is bool and value is int (which comes from Proxmox API),
				// should be converted to bool (as in ".tf" conf).
				case int:
					sValue := strconv.Itoa(value.(int))
					bValue, err := strconv.ParseBool(sValue)
					if err == nil {
						setConfMap[key] = bValue
					}
				// If value is bool, which comes from Terraform conf, add it directly.
				case bool:
					setConfMap[key] = value
				}
			// Anything else will be added as it is.
			default:
				setConfMap[key] = value
			}
			devicesSet.Add(setConfMap)
		}
	}

	return devicesSet
}

// Because default values are not stored in Proxmox, so the API returns only active values.
// So to prevent Terraform doing unnecessary diffs, this function reads default values
// from Terraform itself, and fill empty fields.
func updateDevicesDefaults(
	activeDevicesMap pxapi.QemuDevices,
	configDevicesMap pxapi.QemuDevices,
) pxapi.QemuDevices {

	for deviceID, deviceConf := range configDevicesMap {
		if _, ok := activeDevicesMap[deviceID]; !ok {
			activeDevicesMap[deviceID] = configDevicesMap[deviceID]
		}
		for key, value := range deviceConf {
			if _, ok := activeDevicesMap[deviceID][key]; !ok {
				activeDevicesMap[deviceID][key] = value
			}
		}
	}
	return activeDevicesMap
}

func expandVmQemu(d *schema.ResourceData) pxapi.ConfigQemu {
	vmName := d.Get("name").(string)
	vga := d.Get("vga").(*schema.Set)
	qemuVgaList := vga.List()
	networks := d.Get("network").(*schema.Set)
	qemuNetworks := DevicesSetToMap(networks)
	disks := d.Get("disk").(*schema.Set)
	qemuDisks := DevicesSetToMap(disks)
	serials := d.Get("serial").(*schema.Set)
	qemuSerials := DevicesSetToMap(serials)

	config := pxapi.ConfigQemu{
		Name:         vmName,
		Description:  d.Get("desc").(string),
		Pool:         d.Get("pool").(string),
		Bios:         d.Get("bios").(string),
		Onboot:       d.Get("onboot").(bool),
		Boot:         d.Get("boot").(string),
		BootDisk:     d.Get("bootdisk").(string),
		Agent:        d.Get("agent").(int),
		Memory:       d.Get("memory").(int),
		Balloon:      d.Get("balloon").(int),
		QemuCores:    d.Get("cores").(int),
		QemuSockets:  d.Get("sockets").(int),
		QemuVcpus:    d.Get("vcpus").(int),
		QemuCpu:      d.Get("cpu").(string),
		QemuNuma:     d.Get("numa").(bool),
		Hotplug:      d.Get("hotplug").(string),
		Scsihw:       d.Get("scsihw").(string),
		HaState:      d.Get("hastate").(string),
		QemuOs:       d.Get("qemu_os").(string),
		QemuNetworks: qemuNetworks,
		QemuDisks:    qemuDisks,
		QemuSerials:  qemuSerials,
		// Cloud-init.
		CIuser:       d.Get("ciuser").(string),
		CIpassword:   d.Get("cipassword").(string),
		CIcustom:     d.Get("cicustom").(string),
		Searchdomain: d.Get("searchdomain").(string),
		Nameserver:   d.Get("nameserver").(string),
		Sshkeys:      d.Get("sshkeys").(string),
		Ipconfig0:    d.Get("ipconfig0").(string),
		Ipconfig1:    d.Get("ipconfig1").(string),
		Ipconfig2:    d.Get("ipconfig2").(string),
	}
	if len(qemuVgaList) > 0 {
		config.QemuVga = qemuVgaList[0].(map[string]interface{})
	}

	return config
}
