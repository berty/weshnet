package netmanager

func NewNoopNetManager() *NetManager {
	return NewNetManager(ConnectivityInfo{
		State:   ConnectivityStateOn,
		NetType: ConnectivityNetWifi,
	})
}
