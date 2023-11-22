package main

import (
	"fmt"
)

// NOTE: Update Default Configuration is done
// TODO: Implement Continous Rest Days
// TODO: Debug Weekly Payroll

// Struct for the system's default configuration
type defaultConfig struct {
	daily_salary  float32
	max_reg_hours int
	day_type      []string
	in_time       []string
	out_time      []string
}

// Struct for the regular rates assuming work hours does not exceed 8 hours
type normalMultiplier struct {
	regular         float32
	rest            float32
	nonworking      float32
	nonworking_rest float32
	holiday         float32
	holiday_rest    float32
}

// Struct for the overtime non-night shift rates
type overtimeMultiplier struct {
	normal          float32
	rest            float32
	nonworking      float32
	nonworking_rest float32
	holiday         float32
	holiday_rest    float32
}

// Struct for the overtime night shift rates
type overtimeNightshiftMultiplier struct {
	normal          float32
	rest            float32
	nonworking      float32
	nonworking_rest float32
	holiday         float32
	holiday_rest    float32
}

// Main function
func main() {
	var flag bool = false
	var dc defaultConfig
	var nm normalMultiplier
	var om overtimeMultiplier
	var on overtimeNightshiftMultiplier

	initializeStruct(&dc, &nm, &om, &on)

	for !flag {
		fmt.Println()
		fmt.Println("=========================================")
		fmt.Println("          WEEKLY PAYROLL SYSTEM          ")
		fmt.Println("=========================================")
		fmt.Println()
		fmt.Println("[1] Generate Weekly Payroll")
		fmt.Println("[2] Update Default Configuration")
		fmt.Println("[3] Exit")
		fmt.Println()
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			flag = true
			fmt.Println()
			generateWeeklyPayroll(&dc, &nm, &om, &on)
		case 2:
			fmt.Println()
			updateDefaultConfiguration(&dc)
		case 3:
			flag = true
			fmt.Println()
			fmt.Println("(i) Exiting the program...")
			fmt.Println()
			return
		default:
			fmt.Println()
			fmt.Println("(!) Please enter a valid option")
		}
	}
}

// Generate Weekly Payroll function
func generateWeeklyPayroll(dc *defaultConfig, nm *normalMultiplier, om *overtimeMultiplier, on *overtimeNightshiftMultiplier) {

	var salary float32 = 0
	var weekly_salary float32 = 0

	fmt.Println("-----------------------------------------")
	fmt.Println("         GENERATING WEEKLY PAYROLL       ")
	fmt.Println("-----------------------------------------")
	fmt.Println()

	//var weekly_salary float32 = 0
	var out_times = []string{"", "", "", "", "", "", ""}

	// Get Out Time of each employee for each day of the week
	for i := 0; i < 7; i++ {
		for {
			fmt.Print("Day ", i+1, " Out Time: ")
			fmt.Scanln(&out_times[i])

			if militaryTimeToInt(out_times[i]) < 0 || militaryTimeToInt(out_times[i]) > 24 || len(out_times[i]) != 4 {
				fmt.Println()
				fmt.Println("Please enter a valid time.")
				fmt.Println()
			} else {
				break
			}
		}
	}

	// Get Daily Salary Computation
	for i := 0; i < 7; i++ {
		fmt.Println()
		fmt.Println("-----------------------------------------")
		fmt.Println("                  DAY", i+1, "                  ")
		fmt.Println()

		if militaryTimeToInt(out_times[i]) == militaryTimeToInt(dc.out_time[i]) {
			// If employees are absent in normal days{
			if dc.day_type[i] == "Normal Day" {
				fmt.Println("      Employee is absent on this day")
				fmt.Println()
				continue
			} else if dc.day_type[i] == "Rest Day" {
				fmt.Println("Employee did not come to work")
				// Will still generate the daily rate, etc...
			}
		}

		fmt.Println("Daily Rate:	", dc.daily_salary)
		fmt.Println("IN Time:	", dc.in_time[i])
		fmt.Println("OUT Time:	", out_times[i])
		fmt.Println("Day Type:	", dc.day_type[i])

		work_hours := trackWorkingHours(militaryTimeToInt(out_times[i]), militaryTimeToInt(dc.in_time[i]), dc.max_reg_hours)
		fmt.Println("Night Shift:	", work_hours[0])
		fmt.Println("Overtime:	", work_hours[1])
		fmt.Println("Night Shift OT:	", work_hours[2])

		// Get daily salary
		salary += dc.daily_salary * nmSalary(*nm, dc.day_type[i])

		// Get night shift differential
		if work_hours[0] != 0 {
			salary += float32(work_hours[0]) * (dc.daily_salary / float32(dc.max_reg_hours)) * 1.10
		}

		// Get overtime
		if work_hours[1] != 0 {
			salary += float32(work_hours[1]) * (dc.daily_salary / float32(dc.max_reg_hours)) * omSalary(*om, dc.day_type[i])
		}

		// Get nightshift overtime
		if work_hours[2] != 0 {
			salary += float32(work_hours[2]) * (dc.daily_salary / float32(dc.max_reg_hours)) * onSalary(*on, dc.day_type[i])
		}

		weekly_salary += salary

		// Daily salary
		fmt.Printf("Salary:		%.2f\n\n", salary)
	}
	fmt.Println("-----------------------------------------")
	fmt.Println()
	fmt.Printf("      TOTAL WEEKLY SALARY: %.2f\n", weekly_salary)
	fmt.Println()
	fmt.Println("-----------------------------------------")
}

// Helper function to track work hours
func trackWorkingHours(out_time int, in_time int, max_reg_hours int) []int {

	// Shift the clock-out time if it's on the next day
	if in_time > out_time {
		out_time += 24
	}

	working_hours := out_time - in_time

	// Index mapping for tracked hours:
	// [0] = Night shift differential
	// [1] = Overtime
	// [2] = Night shift overtime

	hours_tracked := []int{0, 0, 0}

	// Loop through the hours worked
	for i := 0; i < working_hours; i++ {
		// Skip the first hour (considered as break)
		if i == 0 {
			in_time++
			continue
		}

		// Categorize hours based on the remaining regular and overtime hours
		if max_reg_hours > 0 { // Check for Night Shift Differential hours
			switch in_time {
			case 22, 23, 24, 1, 2, 3, 4, 5, 6:
				hours_tracked[0]++
			}
			max_reg_hours--

		} else { // Check for Overtime and Night Shift Overtime
			switch in_time {
			case 1, 2, 3, 4, 5, 6, 22, 23, 24:
				hours_tracked[2]++ // Night Shift Overtime
			case 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21:
				hours_tracked[1]++ // Overtime
			}
		}

		// Move to the next hour, reset if it reaches the end of the day (25)
		in_time++
		if in_time == 25 {
			in_time = 1 // Back to the start of the day
		}
	}
	return hours_tracked
}

// Helper functions as salary multipliers (Normal, Overtime, Overtime Nightshift)
func nmSalary(nm normalMultiplier, day_type string) float32 {
	var mult float32

	switch day_type {
	case "Normal Day":
		mult = 1.00
	case "Rest Day":
		mult = nm.rest
	case "Special Non-Working Day":
		mult = nm.nonworking
	case "Special Non-Working and Rest Day":
		mult = nm.nonworking_rest
	case "Regular Holiday":
		mult = nm.holiday
	case "Regular Holiday and Rest Day":
		mult = nm.holiday_rest
	}
	return mult
}

func omSalary(om overtimeMultiplier, day_type string) float32 {
	var mult float32

	switch day_type {
	case "Normal Day":
		mult = om.normal
	case "Rest Day":
		mult = om.rest
	case "Special Non-Working Day":
		mult = om.nonworking
	case "Special Non-Working and Rest Day":
		mult = om.nonworking_rest
	case "Regular Holiday":
		mult = om.holiday
	case "Regular Holiday and Rest Day":
		mult = om.holiday_rest
	}
	return mult
}

func onSalary(on overtimeNightshiftMultiplier, day_type string) float32 {
	var mult float32

	switch day_type {
	case "Normal Day":
		mult = on.normal
	case "Rest Day":
		mult = on.rest
	case "Special Non-Working Day":
		mult = on.nonworking
	case "Special Non-Working and Rest Day":
		mult = on.nonworking_rest
	case "Regular Holiday":
		mult = on.holiday
	case "Regular Holiday and Rest Day":
		mult = on.holiday_rest
	}
	return mult
}

// Update System Configuration Menu
func updateDefaultConfiguration(dc *defaultConfig) {
	for {
		fmt.Println("-----------------------------------------")
		fmt.Println("       UPDATE SYSTEM CONFIGURATION      ")
		fmt.Println("-----------------------------------------")
		fmt.Println()
		fmt.Println("[1] Daily Salary")
		fmt.Println("[2] Max Working Hours Per Day")
		fmt.Println("[3] Day Type")
		fmt.Println("[4] In Time/Out Time")
		fmt.Println("[5] Exit")
		fmt.Println()

		var choice int
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&choice)
		fmt.Println()

		switch choice {
		case 1:
			updateSalary(dc)
		case 2:
			updateWorkHours(dc)
		case 3:
			updateDayType(dc)
		case 4:
			updateInOutTime(dc)
		case 5:
			fmt.Println("(i) Exiting configuration update...")
			return
		default:
			fmt.Println("(!) Please enter a valid option")
		}

		fmt.Println()
		fmt.Println("[1] Continue")
		fmt.Println("[2] Exit")

		var exitChoice int
		fmt.Println()
		fmt.Print("Enter your choice: ")
		fmt.Scanln(&exitChoice)
		fmt.Println()

		switch exitChoice {
		case 1:
			continue
		case 2:
			fmt.Println("(!) Redirecting to Main Menu")
			return
		default:
			fmt.Println("(!) Please enter a valid option")
			fmt.Println()
			continue
		}
	}
}

// Function that updates the daily salary
func updateSalary(dc *defaultConfig) {
	var new_salary float32

	fmt.Println("Current Daily Salary: ", dc.daily_salary)

	fmt.Print("New Daily Salary: ")
	fmt.Scanln(&new_salary)

	dc.daily_salary = new_salary

	fmt.Println()
	fmt.Println("Daily Salary is successfully updated!")
	fmt.Println("    The New Daily Salary is ", dc.daily_salary)
	fmt.Println()
}

// Function that updates the work hours
func updateWorkHours(dc *defaultConfig) {
	var flag bool = false
	var new_max_reg_hours int

	for flag != true {
		fmt.Println("Current Max Working Hours: ", dc.max_reg_hours)

		fmt.Print("New Max Working Hours: ")
		fmt.Scanln(&new_max_reg_hours)

		if new_max_reg_hours < 0 || new_max_reg_hours > 24 {
			fmt.Println()
			fmt.Println("(!) Invalid Working Hours")
			fmt.Println()
		} else {
			flag = true
			dc.max_reg_hours = new_max_reg_hours

			fmt.Println()
			fmt.Println("Max Working Hours is successfully updated!")
			fmt.Println("    The New Max Working Hours is ", dc.max_reg_hours)
			fmt.Println()
		}
	}
}

// Function that updates the day type
func updateDayType(dc *defaultConfig) {
	var chosen_day int
	var new_day_type int

	fmt.Println("Current Day Types:")
	fmt.Println()
	for i := 0; i < len(dc.day_type); i++ {
		fmt.Println("Day", i+1, "Type:", dc.day_type[i])
	}

	for {
		fmt.Println()
		fmt.Print("Choose a Day [1-7]: ")
		fmt.Scanln(&chosen_day)
		fmt.Println()
		if chosen_day < 1 || chosen_day > 7 {
			fmt.Println("Invalid day chosen. Please select a day between 1 and 7.")
		} else {
			break
		}
	}

	fmt.Println("[1] Normal Day")
	fmt.Println("[2] Rest Day")
	fmt.Println("[3] Special Non-Working Day")
	fmt.Println("[4] Special Non-Working  and Rest Day")
	fmt.Println("[5] Regular Holiday")
	fmt.Println("[6] Regular Holiday and Rest Day")
	fmt.Println()

	for {
		fmt.Println()
		fmt.Print("Choose a Day Type: ")
		fmt.Scanln(&new_day_type)
		fmt.Println()
		if new_day_type < 1 || new_day_type > 6 {
			fmt.Println("Please chose among the given options.")
		} else {
			fmt.Println("Day Type is successfully updated!")
			break
		}
	}

	dc.day_type[chosen_day-1] = dayTypeToInt(new_day_type)
}

// Function yhat updates the in and out time
func updateInOutTime(dc *defaultConfig) {
	var chosen_day int
	var new_time string

	fmt.Println("Current Time In and Out Time: ")
	fmt.Println()
	for i := 0; i < len(dc.day_type); i++ {
		fmt.Println("Day", i+1, "Time:", dc.in_time[i])
	}

	for {
		fmt.Println()
		fmt.Print("Choose a Day [1-7]: ")
		fmt.Scanln(&chosen_day)
		fmt.Println()
		if chosen_day < 1 || chosen_day > 7 {
			fmt.Println("Invalid day chosen. Please select a day between 1 and 7.")
		} else {
			break
		}
	}

	for {
		fmt.Print("New In/Out Time: ")
		fmt.Scanln(&new_time)

		if militaryTimeToInt(new_time) < 0 || militaryTimeToInt(new_time) > 24 || len(new_time) != 4 {
			fmt.Println()
			fmt.Println("Please enter a valid time.")
			fmt.Println()
		} else {
			fmt.Println("In and Out Time is successfully updated!")
			break
		}
	}

	dc.in_time[chosen_day-1] = new_time
	dc.out_time[chosen_day-1] = new_time
}

// Helper function to convert day types to int
func dayTypeToInt(input int) string {
	switch input {
	case 1:
		return "Normal Day"
	case 2:
		return "Rest Day"
	case 3:
		return "Special Non-Working Day"
	case 4:
		return "Special Non-Working  and Rest Day"
	case 5:
		return "Regular Holiday"
	case 6:
		return "Regular Holiday and Rest Day"
	default:
		return ""
	}
}

// Helper function to covert military time to int
func militaryTimeToInt(input string) int {

	hours := (int(input[0]-'0') * 10) + int(input[1]-'0')
	minutes := (int(input[2]-'0') * 10) + int(input[3]-'0')

	if hours >= 0 && hours <= 24 && minutes >= 0 && minutes <= 59 {
		return hours
	}
	return -1
}

// Initializes the structs
func initializeStruct(dc *defaultConfig, nm *normalMultiplier, om *overtimeMultiplier, on *overtimeNightshiftMultiplier) {
	dc.daily_salary = 500.00
	dc.max_reg_hours = 8
	dc.day_type = []string{"Normal Day", "Normal Day", "Normal Day", "Normal Day", "Normal Day", "Normal Day", "Normal Day"}
	dc.in_time = []string{"0900", "0900", "0900", "0900", "0900", "0900", "0900"}
	dc.out_time = dc.in_time[:] // copy in time slice since in time and out time of the 7 days is 0900 (system's default configuration)

	nm.regular = 1.00         // REGULAR SALARY
	nm.rest = 1.30            // 130%
	nm.nonworking = 1.30      // 130%
	nm.nonworking_rest = 1.50 // 150%
	nm.holiday = 2.00         // 200%
	nm.holiday_rest = 2.60    // 260%

	nm.regular = 1.25         // 125%
	nm.rest = 1.69            // 169%
	nm.nonworking = 1.69      // 169%
	nm.nonworking_rest = 1.95 // 195%
	nm.holiday = 2.60         // 260%
	nm.holiday_rest = 3.38    // 338%

	nm.regular = 1.375         // 137.5%
	nm.rest = 1.859            // 185.9%
	nm.nonworking = 1.859      // 185.9%
	nm.nonworking_rest = 2.145 // 214.5%
	nm.holiday = 2.860         // 286.0%
	nm.holiday_rest = 3.718    // 371.8%
}
