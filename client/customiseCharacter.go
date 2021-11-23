package client

func customise() {
	player.state = "customise"
}

func (menu *menuSettings) changeModel() {
	menu.playerModel.accessory.model = arrayOfModels["accessory"][menu.chosenModel["accessory"]].model
	menu.playerModel.hair.model = arrayOfModels["hair"][menu.chosenModel["hair"]].model
	menu.playerModel.head.model = arrayOfModels["head"][menu.chosenModel["head"]].model
	menu.playerModel.body.model = arrayOfModels["body"][menu.chosenModel["body"]].model
	menu.playerModel.bottom.model = arrayOfModels["bottom"][menu.chosenModel["bottom"]].model
}
