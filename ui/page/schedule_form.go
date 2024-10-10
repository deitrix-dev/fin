package page

import (
	"fmt"
	"slices"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/murl"
	"github.com/deitrix/fin/ui/api"
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

func repeatOptions(multiplier int) [][2]string {
	if multiplier == 1 {
		return [][2]string{
			{"month", "Month"},
			{"week", "Week"},
			{"day", "Day"},
		}
	}
	return [][2]string{
		{"month", "Months"},
		{"week", "Weeks"},
		{"day", "Days"},
	}
}

func ScheduleForm(
	accounts []string,
	rp fin.RecurringPayment,
	schedule fin.PaymentSchedule,
	in api.ScheduleFormInput,
	index int,
) Node {
	createUpdate := "Create"
	if index > -1 {
		createUpdate = "Update"
	}
	isNewAccount := !slices.Contains(accounts, schedule.AccountID)
	var postURL string
	if index > -1 {
		postURL = fmt.Sprintf("/recurring-payments/%s/schedules/%d", rp.ID, index)
	} else {
		postURL = fmt.Sprintf("/recurring-payments/%s/schedules/new", rp.ID)
	}
	hxFieldRefresh := func(field string) Group {
		return Group{
			hx.Get(murl.Mutate(postURL, murl.AddQuery("formPriority", "true", "sourceField", field))),
			hx.Trigger("change"),
			hx.Include("#scheduleForm"),
			hx.Select("#scheduleForm"),
			hx.Target("#scheduleForm"),
			hx.Swap("outerHTML"),
		}
	}
	return Layout("New schedule",
		Div(Class("grid grid-cols-5 gap-2 overflow-auto"),
			Div(Class("flex overflow-auto items-start col-span-2"),
				Article(Class("flex flex-col flex-1 m-0 bg-white p-8 border-2 border-solid border-gray-300"),
					Div(Class("mb-8"),
						H1(Class("text-center m-0"), Text(createUpdate+" Schedule")),
					),
					Form(ID("scheduleForm"), Action(postURL), Method("post"),
						hx.Get("/api/payments-for-schedule"), hx.Trigger("load, change"),
						hx.Target("#paymentsForSchedule"),
						Div(Class("flex flex-col gap-4"),
							Label(Text("Account"),
								Select(Class("block w-full text-lg p-2 h-12 mt-1"), ID("account"), Name("account"),
									Map(accounts, func(account string) Node {
										return Option(Value(account), Text(account), If(account == schedule.AccountID, Selected()))
									}),
									Option(Value(""), Text("New account"), If(isNewAccount, Selected())),
									hxFieldRefresh("account"),
								),
							),
							If(isNewAccount,
								Label(Text("New account"),
									Input(Type("text"), Class("block w-full text-lg p-2 h-12mt-1"), AutoComplete("off"), ID("newAccount"), Name("newAccount"),
										Value(schedule.AccountID)))),
							Label(Text("Amount"),
								Input(Type("number"), Step("0.01"), Class("block w-full text-lg p-2 h-12 mt-1"), AutoComplete("off"), ID("amount"), Name("amount"), Required(),
									Value(fmt.Sprint(float64(schedule.Amount)/100)))),
							Div(Class("flex gap-2"),
								Label(Class("flex-1"), Text("Start date"),
									Input(Type("date"), Class("block w-full text-lg p-2 h-12 mt-1"), AutoComplete("off"), ID("startDate"), Name("startDate"),
										Iff(schedule.StartDate != nil, func() Node { return Value(schedule.StartDate.Format("2006-01-02")) })),
									hxFieldRefresh("startDate"),
								),
								Label(Class("flex-1"), Text("End date"),
									Input(Type("date"), Class("block w-full text-lg p-2 h-12 mt-1"), AutoComplete("off"), ID("endDate"), Name("endDate"),
										Iff(schedule.EndDate != nil, func() Node { return Value(schedule.EndDate.Format("2006-01-02")) }),
										hxFieldRefresh("endDate"),
									)),
								Label(Class("flex-1"), Text("Payment count"),
									Input(Type("number"), Class("block w-full text-lg p-2 h-12 mt-1"), AutoComplete("off"), ID("count"), Name("count"),
										hxFieldRefresh("count"),
									),
								),
							),
							Div(Class("flex gap-2 items-end"),
								Div(Class("flex flex-col flex-1"),
									Label(Text("Every"), For("every")),
									Input(Type("number"), Class("block w-full text-lg p-2 h-12 mt-1"), AutoComplete("off"), ID("multiplier"), Name("multiplier"),
										Value(fmt.Sprint(schedule.Repeat.Multiplier)),
										hxFieldRefresh("multiplier"),
									),
								),
								Div(Class("flex flex-col flex-1"),
									Select(Class("block w-full text-lg p-2 h-12 mt-1"), ID("repeat"), Name("repeat"), Required(),
										Map(repeatOptions(schedule.Repeat.Multiplier), func(value [2]string) Node {
											return Option(Value(value[0]), Text(value[1]), If(value[0] == string(schedule.Repeat.Every), Selected()))
										}),
										hxFieldRefresh("repeat"),
									),
								),
								If(schedule.Repeat.Every == fin.Month,
									Div(Class("flex flex-col flex-1"),
										Label(Text("Day"), For("dayOfMonth")),
										Input(Type("number"), Class("block w-full text-lg p-2 h-12 mt-1"), AutoComplete("off"), ID("dayOfMonth"), Name("dayOfMonth"),
											Value(fmt.Sprint(schedule.Repeat.Day)),
										),
									),
								),
							),
							If(schedule.Repeat.Multiplier > 1,
								Label(Text("Offset"),
									Input(Type("number"), Class("block w-full text-lg p-2 h-12 mt-1"), AutoComplete("off"), ID("scheduleOffset"), Name("scheduleOffset"),
										Value(fmt.Sprint(schedule.Repeat.Offset))))),

							Button(Class("border-none bg-blue-600 text-white text-lg p-2 h-12 hover:bg-blue-500 cursor-pointer"),
								Type("submit"), Text(createUpdate)),
						),
					),
				),
			),
			Div(Class("flex overflow-auto col-span-3"), ID("paymentsForSchedule")),
		),
	)
}
