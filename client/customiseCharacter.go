package client

func customise() {
	player.state = "customise"
}

func (menu *menuSettings) changeModel() {
	player.model.accessory = arrayOfModels["accessory"][player.chosenModel.Accessory]
	player.model.hair = arrayOfModels["hair"][player.chosenModel.Hair]
	player.model.head = arrayOfModels["head"][player.chosenModel.Head]
	player.model.body = arrayOfModels["body"][player.chosenModel.Body]
	player.model.bottom = arrayOfModels["bottom"][player.chosenModel.Bottom]

}
