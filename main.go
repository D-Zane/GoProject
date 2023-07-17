package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	Id          uint
	Title       string
	Description string
}

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("My ToDoList")
	w.Resize(fyne.NewSize(500, 600))
	w.CenterOnScreen()

	ic, _ := fyne.LoadResourceFromPath("todo.png")
	w.SetIcon(ic)
	var tasks []Task
	var createContent *fyne.Container
	var tasksList *widget.List
	var tasksContent *fyne.Container

	DB, _ := gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	DB.AutoMigrate(&Task{})
	DB.Find(&tasks)

	noTasksLabel := canvas.NewText("Нет Задач", color.Black)

	if len(tasks) != 0 {
		noTasksLabel.Hide()
	}

	newTaskIcon, _ := fyne.LoadResourceFromPath("./icons/icon.png")
	backIcon, _ := fyne.LoadResourceFromPath("./icons/back.png")
	saveIcon, _ := fyne.LoadResourceFromPath("./icons/save.png")
	deleteIcon, _ := fyne.LoadResourceFromPath("./icons/delete.png")
	changeIcon, _ := fyne.LoadResourceFromPath("./icons/change.png")

	tasksBar := container.NewHBox(
		canvas.NewText("Задачи:", color.Black),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", newTaskIcon, func() {
			w.SetContent(createContent)
		}),
	)

	tasksList = widget.NewList(
		func() int {
			return len(tasks)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Defualt")
		},
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(tasks[lii].Title)
		},
	)

	titleEntry := widget.NewEntry()
	titleEntry.SetPlaceHolder("Название задачу")

	descriptionEntry := widget.NewMultiLineEntry()
	descriptionEntry.SetPlaceHolder("Описание")

	tasksList.OnSelected = func(id widget.ListItemID) {
		detailsBar := container.NewHBox(
			canvas.NewText(
				fmt.Sprintf(
					"Details about \"%s\"",
					tasks[id].Title,
				),
				color.Black,
			),
			layout.NewSpacer(),
			widget.NewButtonWithIcon("", backIcon, func() {
				w.SetContent(tasksContent)
				tasksList.Unselect(id)
			}),
		)

		tasktTitle := widget.NewLabel(tasks[id].Title)
		tasktTitle.TextStyle = fyne.TextStyle{Bold: true}

		taskDescriptiont := widget.NewLabel(tasks[id].Description)
		taskDescriptiont.TextStyle = fyne.TextStyle{Italic: true}
		taskDescriptiont.Wrapping = fyne.TextWrapBreak

		buttonsBox := container.NewHBox(
			// Удаление
			widget.NewButtonWithIcon(
				"",
				deleteIcon,

				func() {
					dialog.ShowConfirm(
						"Удалить задачу",
						fmt.Sprintf(
							"Вы уверены, что хотите удалить задачу\"%s\"?",
							tasks[id].Title,
						),

						func(b bool) {
							if b {
								DB.Delete(&Task{}, "Id = ?", tasks[id].Id)
								DB.Find(&tasks)

								if len(tasks) == 0 {
									noTasksLabel.Show()
								} else {
									noTasksLabel.Hide()
								}
							}

							w.SetContent(tasksContent)
							tasksList.UnselectAll()
						},
						w,
					)
				},
			),

			//Изменение
			widget.NewButtonWithIcon(
				"",
				changeIcon,

				func() {
					changeBar := container.NewHBox(
						canvas.NewText(
							fmt.Sprintf(
								"Изменение \"%s\"",
								tasks[id].Title,
							),
							color.Black,
						),
						layout.NewSpacer(),
						widget.NewButtonWithIcon("", backIcon, func() {
							w.SetContent(tasksContent)
							tasksList.Unselect(id)
						}),
					)

					changeTitle := widget.NewEntry()
					changeTitle.SetText(tasks[id].Title)

					changeDescription := widget.NewEntry()
					changeDescription.SetText(tasks[id].Description)
					changeDescription.Wrapping = fyne.TextWrapBreak

					changeButton := widget.NewButtonWithIcon(
						"Сохранить задачу",
						saveIcon,
						func() {
							DB.Find(&Task{}, "Id = ?", tasks[id].Id).Updates(Task{
								Title:       changeTitle.Text,
								Description: changeDescription.Text,
							},
							)

							DB.Find(&tasks)

							w.SetContent(tasksContent)
							tasksList.UnselectAll()
						},
					)

					changeContent := container.NewVBox(
						changeBar,
						canvas.NewLine(color.Black),

						changeTitle,
						changeDescription,
						changeButton,
					)

					w.SetContent(changeContent)
				},
			),
		)

		detailsVBox := container.NewVBox(
			detailsBar,
			canvas.NewLine(color.Black),

			tasktTitle,
			taskDescriptiont,
			buttonsBox,
		)

		w.SetContent(detailsVBox)
	}

	tasksScroll := container.NewScroll(tasksList)
	tasksScroll.SetMinSize(fyne.NewSize(500, 500))

	tasksContent = container.NewVBox(
		tasksBar,
		canvas.NewLine(color.Black),
		noTasksLabel,
		tasksScroll,
	)

	saveTaskButton := widget.NewButtonWithIcon("Сохранить задачу", saveIcon, func() {
		task := Task{
			Title:       titleEntry.Text,
			Description: descriptionEntry.Text,
		}

		DB.Create(&task)
		DB.Find(&tasks)

		//Обновление строк
		titleEntry.Text = ""
		titleEntry.Refresh()
		//Обновление строк
		descriptionEntry.Text = ""
		descriptionEntry.Refresh()

		w.SetContent(tasksContent)

		tasksList.UnselectAll()

		if len(tasks) == 0 {
			noTasksLabel.Show()
		} else {
			noTasksLabel.Hide()
		}
	})

	createBar := container.NewHBox(
		canvas.NewText("Задачи:", color.Black),
		layout.NewSpacer(),
		widget.NewButtonWithIcon("", backIcon, func() {

			//Обновление строк
			titleEntry.Text = ""
			titleEntry.Refresh()
			//Обновление строк
			descriptionEntry.Text = ""
			descriptionEntry.Refresh()

			w.SetContent(tasksContent)
			tasksList.UnselectAll()
		}),
	)

	createContent = container.NewVBox(
		createBar,
		canvas.NewLine(color.Black),

		container.NewVBox(
			titleEntry,
			descriptionEntry,
			saveTaskButton,
		),
	)

	w.SetContent(tasksContent)

	w.Show()
	a.Run()
}
