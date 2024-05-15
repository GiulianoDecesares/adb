package fastboot

type Device struct {
	serialNo    string
	fastbootCli *Fastboot
}

func NewDevice(id string, fastboot *Fastboot) *Device {
	return &Device{serialNo: id, fastbootCli: fastboot}
}

func (device *Device) GetSerial() string { return device.serialNo }

func (device *Device) Boot() error {
	_, err := device.Run("continue")
	return err
}

func (device *Device) Reboot() error {
	_, err := device.Run("reboot")
	return err
}

func (device *Device) Run(rawCommand ...string) (string, error) {
	var arguments []string

	arguments = append(arguments, "-s", device.serialNo)
	arguments = append(arguments, rawCommand...)

	return device.fastbootCli.Run(arguments...)
}
