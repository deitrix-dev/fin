package api

import (
	"net/http"
	"time"

	"github.com/deitrix/fin"
	"github.com/deitrix/fin/pkg/form"
	"github.com/deitrix/fin/pkg/iterx"
	"github.com/deitrix/fin/pkg/pointer"
	"github.com/deitrix/fin/ui/components"
)

type ScheduleFormInput struct {
	Account        string
	NewAccount     string
	Amount         float64
	StartDate      time.Time
	EndDate        time.Time
	Repeat         string
	DayOfMonth     int
	DayOfWeek      string
	ScheduleOffset int
	Multiplier     int
	Offset         int
	OOB            bool
	FormPriority   bool
}

func (in ScheduleFormInput) Schedule() fin.PaymentSchedule {
	schedule := fin.PaymentSchedule{
		Repeat: fin.Repeat{
			Every:      fin.Step(in.Repeat),
			Multiplier: in.Multiplier,
			Offset:     in.ScheduleOffset,
		},
		Amount: int(in.Amount * 100),
	}

	switch in.Repeat {
	case "month":
		schedule.Repeat.Day = in.DayOfMonth
	case "week":
		schedule.Repeat.Weekday = fin.Weekday(in.DayOfWeek)
	}

	if !in.StartDate.IsZero() && in.StartDate.Unix() > 0 {
		schedule.StartDate = &in.StartDate
	}
	if !in.EndDate.IsZero() && in.EndDate.Unix() > 0 {
		schedule.EndDate = &in.EndDate
	}

	if in.Account != "" {
		schedule.AccountID = in.Account
	} else {
		schedule.AccountID = in.NewAccount
	}

	return schedule
}

func ScheduleForm(in *ScheduleFormInput) form.Fields {
	return form.Fields{
		"account":        form.String(&in.Account),
		"newAccount":     form.String(&in.NewAccount),
		"amount":         form.Float(&in.Amount),
		"startDate":      form.Time(&in.StartDate, "2006-01-02"),
		"endDate":        form.Time(&in.EndDate, "2006-01-02"),
		"repeat":         form.String(&in.Repeat),
		"dayOfMonth":     form.Int(&in.DayOfMonth),
		"dayOfWeek":      form.String(&in.DayOfWeek),
		"scheduleOffset": form.Int(&in.ScheduleOffset),
		"multiplier":     form.Int(&in.Multiplier),
		"offset":         form.Int(&in.Offset),
		"oob":            form.Bool(&in.OOB),
		"formPriority":   form.Bool(&in.FormPriority),
	}
}

func PaymentsForSchedule(w http.ResponseWriter, r *http.Request) {
	var input ScheduleFormInput
	formFields := ScheduleForm(&input)
	if err := form.Decode(r.URL.Query(), formFields); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	from := time.Now()
	if !input.StartDate.IsZero() && input.StartDate.Unix() > 0 {
		from = input.StartDate
	}

	schedule := input.Schedule()
	payments := iterx.CollectN(iterx.Skip(schedule.PaymentsSince(from), input.Offset), 26)
	var nextPage *int
	if len(payments) == 26 {
		nextPage = pointer.To(input.Offset + 25)
		payments = payments[:25]
	}
	fetchURL := "/api/payments-for-schedule?oob=true&" + form.Encode(formFields).Encode()
	renderGomp(w, r, components.Payments(components.PaymentsInputs{
		Header:   "Payments Preview",
		Payments: payments,
		FetchURL: fetchURL,
		NextPage: nextPage,
		OOB:      input.OOB,
	}))
}
