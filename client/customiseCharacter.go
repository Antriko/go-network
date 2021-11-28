package client

func customise() {
	player.state = "customise"
}

func (menu *menuSettings) changeModel() {
	player.model.accessory = arrayOfModels["accessory"][player.UserModelSelection.Accessory]
	player.model.hair = arrayOfModels["hair"][player.UserModelSelection.Hair]
	player.model.head = arrayOfModels["head"][player.UserModelSelection.Head]
	player.model.body = arrayOfModels["body"][player.UserModelSelection.Body]
	player.model.bottom = arrayOfModels["bottom"][player.UserModelSelection.Bottom]

}
