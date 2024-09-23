package components

import (
	"fmt"

	"github.com/deitrix/fin"
	. "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	. "github.com/maragudk/gomponents/html"
)

var weekdayOptions = [][2]string{
	{"monday", "Monday"},
	{"tuesday", "Tuesday"},
	{"wednesday", "Wednesday"},
	{"thursday", "Thursday"},
	{"friday", "Friday"},
	{"saturday", "Saturday"},
	{"sunday", "Sunday"},
}

var repeatOptions = [][2]string{
	{"month", "Monthly"},
	{"week", "Weekly"},
	{"day", "Daily"},
}

func ScheduleForm(
	accounts []string,
	rp fin.RecurringPayment,
	schedule fin.PaymentSchedule,
	index int,
) Node {
	var postURL string
	if index > -1 {
		postURL = fmt.Sprintf("/recurring-payments/%s/schedules/%d", rp.ID, index)
	} else {
		postURL = fmt.Sprintf("/recurring-payments/%s/schedules/new", rp.ID)
	}
	refreshFormURL := postURL + "?formPriority=true"
	return Div(Class("grid grid-cols-3 gap-2 overflow-auto"),
		Div(Class("flex overflow-auto items-start"),
			Article(Class("flex flex-col flex-1 m-0 bg-white p-8 border-2 border-solid border-gray-300"),
				Div(Class("mb-8"),
					H1(Class("text-center m-0"), Text("Create schedule")),
				),
				Form(ID("scheduleForm"), Action(postURL), Method("post"),
					hx.Get("/api/payments-for-schedule"), hx.Trigger("load, change"),
					hx.Target("#paymentsForSchedule"),
					Div(Class("flex flex-col gap-4"),
						Label(Text("Account"),
							Select(Class("block w-full text-lg px-2 py-3.5 mt-1"), ID("account"), Name("account"), Required(),
								Group(Map(accounts, func(account string) Node {
									return Option(Value(account), Text(account), If(account == schedule.AccountID, Selected()))
								})),
							),
						),
						Label(Text("Amount"),
							Input(Type("number"), Step("0.01"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("amount"), Name("amount"), Required(),
								Value(fmt.Sprint(float64(schedule.Amount)/100)))),
						Label(Text("Start date"),
							Input(Type("date"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("startDate"), Name("startDate"),
								Iff(schedule.StartDate != nil, func() Node { return Value(schedule.StartDate.Format("2006-01-02")) }))),
						Label(Text("End date"),
							Input(Type("date"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("endDate"), Name("endDate"),
								Iff(schedule.EndDate != nil, func() Node { return Value(schedule.EndDate.Format("2006-01-02")) }))),
						Label(Text("Repeat"),
							Select(Class("block w-full text-lg px-2 py-3.5 mt-1"), ID("repeat"), Name("repeat"), Required(),
								Group(Map(repeatOptions, func(value [2]string) Node {
									return Option(Value(value[0]), Text(value[1]), If(value[0] == string(schedule.Repeat.Every), Selected()))
								})),
								hx.Get(refreshFormURL), hx.Trigger("change"), hx.Include("#scheduleForm"),
								hx.Select("#scheduleForm"), hx.Target("#scheduleForm"), hx.Swap("outerHTML"),
							)),
						If(schedule.Repeat.Every == fin.Month,
							Label(Text("Day of month"),
								Input(Type("number"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("dayOfMonth"), Name("dayOfMonth"),
									Value(fmt.Sprint(schedule.Repeat.Day))))),
						If(schedule.Repeat.Every == fin.Week,
							Label(Text("Day of week"),
								Select(Class("block w-full text-lg px-2 py-3.5 mt-1"), ID("dayOfWeek"), Name("dayOfWeek"),
									Group(Map(weekdayOptions, func(value [2]string) Node {
										return Option(Value(value[0]), Text(value[1]), If(value[0] == string(schedule.Repeat.Weekday), Selected()))
									}))))),
						Label(Text("Multiplier"),
							Input(Type("number"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("multiplier"), Name("multiplier"),
								Value(fmt.Sprint(schedule.Repeat.Multiplier)),
								hx.Get(refreshFormURL), hx.Trigger("change"), hx.Include("#scheduleForm"),
								hx.Select("#scheduleForm"), hx.Target("#scheduleForm"), hx.Swap("outerHTML"),
							)),
						If(schedule.Repeat.Multiplier > 1,
							Label(Text("Offset"),
								Input(Type("number"), Class("block w-full text-lg p-2 mt-1"), AutoComplete("off"), ID("scheduleOffset"), Name("scheduleOffset"),
									Value(fmt.Sprint(schedule.Repeat.Offset))))),

						Button(Class("border-none bg-blue-600 text-white text-lg py-2 hover:bg-blue-500 cursor-pointer"),
							Type("submit"), If(index == -1, Text("Create")), If(index != -1, Text("Update"))),
					),
				),
			),
		),
		Div(Class("flex overflow-auto col-span-2"), ID("paymentsForSchedule")),
	)
}
