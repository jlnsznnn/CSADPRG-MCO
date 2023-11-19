package main

import (
	"fmt"
)

// NOTE: Update Default Configuration is done
// TODO: Implement Continous Rest Days
// TODO: Generate Weekly Payroll

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

func generateWeeklyPayroll(dc *defaultConfig, nm *normalMultiplier, om *overtimeMultiplier, on *overtimeNightshiftMultiplier) {
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
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
		fmt.Println("                  DAY", i+1, "                  ")
		fmt.Println()

		if militaryTimeToInt(out_times[i]) == militaryTimeToInt(dc.out_time[i]) {
			// If employees are absent in normal days{
			if dc.day_type[i] == "Normal Day" {
				fmt.Println("      Employee is absent on this day")
				fmt.Println()
				continue
				// If emplyees did not work in rest days
			} else if dc.day_type[i] == "Rest Day" {
				fmt.Println("Employee did not come to work (Rest Day)")
			}
		}

		fmt.Println("Daily Rate:	", dc.daily_salary)
		fmt.Println("IN Time:	", dc.in_time)
		fmt.Println("OUT Time:	", dc.out_time)
		fmt.Println("Day Type:	", dc.day_type)
	}
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~")
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

	nm.regular = 1.30         // NOTE: need to clarify to sir
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
