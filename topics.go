package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "victron"

type mqttObserver func(componentType string, componentId string, value float64)

var labels = []string{"component_type", "component_id"}

func gaugeObserver(opts prometheus.GaugeOpts) mqttObserver {
	opts.Namespace = namespace
	gauge := prometheus.NewGaugeVec(opts, labels)
	prometheus.MustRegister(gauge)
	return func(componentType string, componentId string, value float64) {
		gauge.WithLabelValues(componentType, componentId).Set(value)
	}
}

func counterObserver(opts prometheus.CounterOpts) mqttObserver {
	opts.Namespace = namespace
	counter := prometheus.NewCounterVec(opts, labels)
	prometheus.MustRegister(counter)
	var prevValue float64
	first := true
	return func(componentType string, componentId string, value float64) {
		if first {
			prevValue = value
			first = false
			return
		}

		if prevValue <= value {
			counter.WithLabelValues(componentType, componentId).Add(value - prevValue)
		}
		prevValue = value
	}
}

func alarm(alarmType string) mqttObserver {
	gauge := prometheus.GaugeOpts{
		Name:        "alarm",
		Help:        "0=OK; 1=Warning; 2=Alarm",
		ConstLabels: prometheus.Labels{"alarm_type": alarmType},
	}

	return gaugeObserver(gauge)
}

func phaseAlarm(phase string, alarmType string) mqttObserver {
	gauge := prometheus.GaugeOpts{
		Name:        "phase_alarm",
		Help:        "0=OK; 1=Warning; 2=Alarm",
		ConstLabels: prometheus.Labels{"phase": phase, "alarm_type": alarmType},
	}

	return gaugeObserver(gauge)
}

var suffixTopicMap = map[string]mqttObserver{
	"Ac/ActiveIn/Source": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_activein_source",
			Help: "The active AC-In source of the multi",
		}),
	"Ac/Consumption/NumberOfPhases": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_consumption_number_of_phases",
			Help: "",
		}),
	"Ac/Consumption/L1/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_phase_power_watts",
			Help:        "Total of ConsumptionOnInput & ConsumptionOnOutput",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/Consumption/L2/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_phase_power_watts",
			Help:        "Total of ConsumptionOnInput & ConsumptionOnOutput",
			ConstLabels: prometheus.Labels{"phase": "2"},
		}),
	"Ac/Consumption/L3/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_phase_power_watts",
			Help:        "Total of ConsumptionOnInput & ConsumptionOnOutput",
			ConstLabels: prometheus.Labels{"phase": "3"},
		}),
	"Ac/ConsumptionOnInput/NumberOfPhases": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_consumption_on_input_number_of_phases",
			Help: "",
		}),
	"Ac/ConsumptionOnInput/L1/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_input_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/ConsumptionOnInput/L2/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_input_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "2"},
		}),
	"Ac/ConsumptionOnInput/L3/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_input_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "3"},
		}),
	"Ac/ConsumptionOnOutput/NumberOfPhases": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_consumption_on_output_number_of_phases",
			Help: "",
		}),
	"Ac/ConsumptionOnOutput/L1/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_output_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/ConsumptionOnOutput/L2/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_output_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "2"},
		}),
	"Ac/ConsumptionOnOutput/L3/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_output_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "3"},
		}),
	"Dc/Battery/Alarms/CircuitBreakerTripped": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_battery_alarms_circuit_breaker_tripped",
			Help: "",
		}),
	"Dc/Battery/ConsumedAmphours": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_battery_consumed_amphours",
			Help: "Ah",
		}),
	"Dc/Battery/Current": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_battery_current",
			Help: "",
		}),
	"Dc/Battery/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_battery_power_watts",
			Help: "",
		}),
	"Dc/Battery/Soc": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_battery_state_of_charge",
			Help: "",
		}),
	"Dc/Battery/State": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_battery_state",
			Help: "",
		}),
	"Dc/Battery/TimeToGo": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_battery_time_to_go_seconds",
			Help: "",
		}),
	"Dc/Battery/Voltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_battery_voltage_volts",
			Help: "",
		}),
	"Dc/Charger/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_charger_power_watts",
			Help: "",
		}),
	"Dc/Pv/Current": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_pv_current_amps",
			Help: "",
		}),
	"Dc/Pv/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_pv_power_watts",
			Help: "",
		}),
	"Dc/System/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_system_power_watts",
			Help: "",
		}),
	"Dc/Vebus/Current": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_vebus_current_amps",
			Help: "",
		}),
	"Dc/Vebus/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_vebus_power_watts",
			Help: "",
		}),
	"Buzzer/State": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "buzzer_state",
			Help: "",
		}),
	"Relay/0/State": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "relay_state",
			Help:        "",
			ConstLabels: prometheus.Labels{"relay": "0"},
		}),
	"Relay/1/State": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "relay_state",
			Help:        "",
			ConstLabels: prometheus.Labels{"relay": "1"},
		}),
	"SystemState/State": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "system_state",
			Help: "",
		}),
	"Timers/TimeOnGrid": counterObserver(
		prometheus.CounterOpts{
			Name: "time_on_grid_seconds_total",
			Help: "Time spent on grid",
		}),
	"Timers/TimeOnGenerator": counterObserver(
		prometheus.CounterOpts{
			Name: "time_on_generator_seconds_total",
			Help: "Time spent on generator",
		}),
	"Timers/TimeOnInverter": counterObserver(
		prometheus.CounterOpts{
			Name: "time_on_inverter_seconds_total",
			Help: "Time spent on inverter",
		}),
	"Timers/TimeOff": counterObserver(
		prometheus.CounterOpts{
			Name: "time_off_seconds_total",
			Help: "Time spent off",
		}),
	"Settings/CGwacs/AcPowerSetPoint": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_ac_power_set_point",
			Help: "User setting: Grid set-point",
		}),
	"Settings/CGwacs/BatteryLife/DischargedSoc": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_discharged_state_of_charge",
			Help: "Deprecated",
		}),
	"Settings/CGwacs/BatteryLife/DischargedTime": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_dischanged_time",
			Help: "Internal",
		}),
	"Settings/CGwacs/BatteryLife/Flags": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_flags",
			Help: "Internal",
		}),
	"Settings/CGwacs/BatteryLife/MinimumSocLimit": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_minimum_state_of_charge_limit",
			Help: "User setting: Minimum Discharge SOC",
		}),
	"Settings/CGwacs/BatteryLife/SocLimit": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_state_of_charge_limit",
			Help: "Output of the BatteryLife algorithm (read only)",
		}),
	"Settings/CGwacs/BatteryLife/State": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_state",
			Help: "ESS state (read & write, see below)",
		}),
	"Settings/CGwacs/Hub4Mode": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_hub4_mode",
			Help: "ESS mode (read & write, see below)",
		}),
	"Settings/CGwacs/MaxChargePercentage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_max_charge_percentage",
			Help: "Deprecated",
		}),
	"Settings/CGwacs/MaxChargePower": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_max_charge_power_watts",
			Help: "User setting: Max Charge Power",
		}),
	"Settings/CGwacs/MaxDischargePercentage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_max_discharge_percentage",
			Help: "Deprecated",
		}),
	"Settings/CGwacs/MaxDischargePower": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_max_discharge_power_watts",
			Help: "User setting: Max Inverter Power",
		}),
	"Settings/CGwacs/OvervoltageFeedIn": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_overvoltage_feed_in",
			Help: "User setting: Feed-in excess solar charger power (yes/no)",
		}),
	"Settings/CGwacs/PreventFeedback": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_prevent_feedback",
			Help: "User setting: PV Inverter Zero Feed-in (on/off)",
		}),
	"Settings/CGwacs/RunWithoutGridMeter": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_run_without_grid_meter",
			Help: "User setting: Grid meter installed (on/off)",
		}),
	/** com.victronenergy.vebus */
	"Ac/ActiveIn/L1/F": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_active_input_phase__freq_hz",
			Help:        "Frequency",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/ActiveIn/L1/I": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_active_input_phase_current_amps",
			Help:        "Current",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/ActiveIn/L1/P": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_active_input_phase_power_watts",
			Help:        "Real power",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/ActiveIn/L1/V": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_active_input_phase_voltage_volts",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/ActiveIn/P": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_active_input_power_watts",
			Help: "Total power",
		}),
	"Ac/ActiveIn/Connected": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_active_input_connected",
			Help: "0 when inverting, 1 when connected to an AC in.",
		}),
	"Ac/ActiveIn/ActiveInput": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_active_input_active_input",
			Help: "Active input: 0 = ACin-1, 1 = ACin-2, 240 is none (inverting).",
		}),
	"Ac/In/1/CurrentLimit": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_input_current_limit",
			Help:        "",
			ConstLabels: prometheus.Labels{"input": "1"},
		}),
	"Ac/In/1/CurrentLimitIsAdjustable": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_input_current_limit_is_adjustable",
			Help:        "",
			ConstLabels: prometheus.Labels{"input": "1"},
		}),
	"Ac/In/2/CurrentLimit": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_input_current_limit",
			Help:        "",
			ConstLabels: prometheus.Labels{"input": "2"},
		}),
	"Ac/In/2/CurrentLimitIsAdjustable": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_input_current_limit_is_adjustable",
			Help:        "",
			ConstLabels: prometheus.Labels{"input": "2"},
		}),
	"Ac/PowerMeasurementType": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_power_measurement_type",
			Help: "Indicates the type of power measurement used by the system.",
		}),
	"Alarms/LowBattery":         alarm("LowBattery"),
	"Alarms/PhaseRotation":      alarm("PhaseRotation"),
	"Alarms/Ripple":             alarm("Ripple"),
	"Alarms/TemperatureSensor":  alarm("TemperatureSensor"),
	"Alarms/L1/HighTemperature": phaseAlarm("1", "HighTemperature"),
	"Alarms/L1/LowBattery":      phaseAlarm("1", "LowBattery"),
	"Alarms/L1/Overload":        phaseAlarm("1", "Overload"),
	"Alarms/L1/Ripple":          phaseAlarm("1", "Ripple"),
	"Alarms/L2/HighTemperature": phaseAlarm("2", "HighTemperature"),
	"Alarms/L2/LowBattery":      phaseAlarm("2", "LowBattery"),
	"Alarms/L2/Overload":        phaseAlarm("2", "Overload"),
	"Alarms/L2/Ripple":          phaseAlarm("2", "Ripple"),
	"Alarms/L3/HighTemperature": phaseAlarm("3", "HighTemperature"),
	"Alarms/L3/LowBattery":      phaseAlarm("3", "LowBattery"),
	"Alarms/L3/Overload":        phaseAlarm("3", "Overload"),
	"Alarms/L3/Ripple":          phaseAlarm("3", "Ripple"),
	"Dc/0/Voltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_voltage_volts",
			Help:        "V DC",
			ConstLabels: prometheus.Labels{"n": "0"},
		}),
	"Dc/0/Current": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_current_amps",
			Help:        "A DC",
			ConstLabels: prometheus.Labels{"n": "0"},
		}),
	"Dc/0/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_power_watts",
			Help:        "",
			ConstLabels: prometheus.Labels{"n": "0"},
		}),
	"Dc/0/Temperature": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_temperature_celcius",
			Help:        "°C - Battery temperature",
			ConstLabels: prometheus.Labels{"n": "0"},
		}),
	"Mode": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "mode",
			Help: "Position of the switch. 1=Charger Only;2=Inverter Only;3=On;4=Off",
		}),
	"ModeIsAdjustable": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "mode_is_adjustable",
			Help: "",
		}),
	"VebusChargeState": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "vebus_charge_state",
			Help: "1. Bulk, 2. Absorption, 3. Float, 4. Storage, 5. Repeat absorption, 6. Forced absorption, 7. Equalise, 8. Bulk stopped",
		}),
	"VebusSetChargeState": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "vebus_set_charge_state",
			Help: "1. Force to Equalise. 2. Force to Absorption, for maximum absorption time. 3. Force to Float, for 24 hours.",
		}),
	"Leds/Mains": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "led_mains",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		}),
	"Leds/Bulk": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "led_bulk",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		}),
	"Leds/Absorption": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "led_absoption",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		}),
	"Leds/Float": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "led_float",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		}),
	"Leds/Inverter": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "led_inverter",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		}),
	"Leds/Overload": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "led_overload",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		}),
	"Leds/LowBattery": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "led_low_battery",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		}),
	"Leds/Temperature": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "led_temperature",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		}),
	/* com.victronenergy.inverter */
	"Alarms/LowVoltage":       alarm("LowVoltage"),
	"Alarms/HighVoltage":      alarm("HighVoltage"),
	"Alarms/LowTemperature":   alarm("LowTemperature"),
	"Alarms/HighTemperature":  alarm("HighTemperature"),
	"Alarms/Overload":         alarm("Overload"),
	"Alarms/LowVoltageAcOut":  alarm("LowVoltageAcOut"),
	"Alarms/HighVoltageAcOut": alarm("HighVoltageAcOut"),
	"Ac/Out/P": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_output_power_watts",
			Help: "AC Output power watts",
		}),
	"Ac/Out/L1/V": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_output_phase_volts",
			Help:        "AC Output voltage",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/Out/L1/I": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_output_phase_current_amps",
			Help:        "AC Output current",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/Out/L1/F": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_output_phase_freq_hz",
			Help:        "AC Output frequency Hertz",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/Out/L1/P": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_output_phase_power_watts",
			Help:        "Not used on vedirect inverters ",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"State": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "state",
			Help: "",
		}),
	"Dc/0/MidVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_midvoltage_volts",
			Help:        "V DC Mid voltage (BMV-702 configured to read midpoint voltage only)",
			ConstLabels: prometheus.Labels{"n": "0"},
		}),
	"Dc/0/MidVoltageDeviation": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_midvoltage_deviation_percent",
			Help:        "Percentage deviation",
			ConstLabels: prometheus.Labels{"n": "0"},
		}),
	"Dc/1/Voltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_voltage_volts",
			Help:        "V DC",
			ConstLabels: prometheus.Labels{"n": "1"},
		}),
	"ConsumedAmphours": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "consumed_amphours",
			Help: "Ah",
		}),
	"Soc": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "state_of_charge",
			Help: "0 to 100 % (BMV, BYD, Lynx BMS)",
		}),
	"TimeToGo": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "time_to_go_seconds",
			Help: "Time to in seconds (BMV SOC relay/discharge floor value, Lynx BMS).  Max value 864,000 when battery is not discharging.",
		}),
	"Info/MaxChargeCurrent": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "max_charge_current_amps",
			Help: "Charge Current Limit aka CCL  (BYD, Lynx BMS and FreedomWon)",
		}),
	"Info/MaxDischargeCurrent": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "max_discharge_current_amps",
			Help: "Discharge Current Limit aka DCL (BYD, Lynx BMS and FreedomWon)",
		}),
	"Info/MaxChargeVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "max_charge_voltage_volts",
			Help: "Maximum voltage to charge to (BYD, Lynx BMS and FreedomWon)",
		}),
	"Info/BatteryLowVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "battery_low_voltage",
			Help: "Note that Low Voltage is ignored by the system (BYD, Lynx BMS and FreedomWon)",
		}),
	"Ac/Alarms/GridLost":           alarm("GridLost"),
	"Alarms/Alarm":                 alarm("Alarm"),
	"Alarms/LowStarterVoltage":     alarm("LowStarterVoltage"),
	"Alarms/HighStarterVoltage":    alarm("HighStarterVoltage"),
	"Alarms/LowSoc":                alarm("LowSoc"),
	"Alarms/HighChargeCurrent":     alarm("HighChargeCurrent"),
	"Alarms/HighDischargeCurrent":  alarm("HighDischargeCurrent"),
	"Alarms/CellImbalance":         alarm("CellImbalance"),
	"Alarms/InternalFailure":       alarm("InternalFailure"),
	"Alarms/HighChargeTemperature": alarm("HighChargeTemperature"),
	"Alarms/LowChargeTemperature":  alarm("LowChargeTemperature"),
	"Alarms/LowCellVoltage":        alarm("LowCellVoltage"),
	"Alarms/MidVoltage":            alarm("MidVoltage"),
	"Settings/HasTemperature": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_has_temperature",
			Help: "",
		}),
	"Settings/HasStarterVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_has_starter_voltage",
			Help: "",
		}),
	"Settings/HasMidVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "settings_has_mid_voltage",
			Help: "",
		}),
	"History/DeepestDischarge": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_deepest_discharge",
			Help: "",
		}),
	"History/LastDischarge": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_last_discharge",
			Help: "",
		}),
	"History/AverageDischarge": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_avg_discharge",
			Help: "",
		}),
	"History/ChargeCycles": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_charge_cycles",
			Help: "",
		}),
	"History/FullDischarges": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_full_discharges",
			Help: "",
		}),
	"History/TotalAhDrawn": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_total_drawn_amphours",
			Help: "",
		}),
	"History/MinimumVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_min_voltage_volts",
			Help: "",
		}),
	"History/MaximumVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_max_voltage_volts",
			Help: "",
		}),
	"History/TimeSinceLastFullCharge": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_time_since_full_charge_seconds",
			Help: "",
		}),
	"History/AutomaticSyncs": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_automatic_syncs",
			Help: "",
		}),
	"History/LowVoltageAlarms": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_low_voltage_alarms",
			Help: "",
		}),
	"History/HighVoltageAlarms": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_high_voltage_alarms",
			Help: "",
		}),
	"History/LowStarterVoltageAlarms": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_low_starter_voltage_alarms",
			Help: "",
		}),
	"History/HighStarterVoltageAlarms": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_high_starter_voltage_alarms",
			Help: "",
		}),
	"History/MinimumStarterVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_min_starter_voltage",
			Help: "",
		}),
	"History/MaximumStarterVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_max_starter_voltage",
			Help: "",
		}),
	"History/DischargedEnergy": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_discharge_energy_kwh",
			Help: "",
		}),
	"History/ChargedEnergy": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_charged_energy_kwh",
			Help: "",
		}),
	"ErrorCode": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "error_code",
			Help: "",
		}),
	"SystemSwitch": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "system_switch",
			Help: "",
		}),
	"Balancing": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "balancing",
			Help: "",
		}),
	"System/NrOfBatteries": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "system_battery_count",
			Help: "",
		}),
	"System/BatteriesParallel": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "system_batteries_parallel_count",
			Help: "",
		}),
	"System/BatteriesSeries": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "system_batteries_series_count",
			Help: "",
		}),
	"System/NrOfCellsPerBattery": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "system_cells_per_battery_count",
			Help: "",
		}),
	"System/MinCellVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "system_min_cell_voltage_volts",
			Help: "",
		}),
	"System/MaxCellVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "system_max_cell_voltage_volts",
			Help: "",
		}),
	"Diagnostics/ShutDownsDueError": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "diagnostics_shutdowns_due_to_error_count",
			Help: "",
		}),
	"Diagnostics/LastErrors/1/Error": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "diagnostics_last_error",
			Help:        "",
			ConstLabels: prometheus.Labels{"e": "1"},
		}),
	"Diagnostics/LastErrors/2/Error": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "diagnostics_last_error",
			Help:        "",
			ConstLabels: prometheus.Labels{"e": "2"},
		}),
	"Diagnostics/LastErrors/3/Error": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "diagnostics_last_error",
			Help:        "",
			ConstLabels: prometheus.Labels{"e": "3"},
		}),
	"Diagnostics/LastErrors/4/Error": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "diagnostics_last_error",
			Help:        "",
			ConstLabels: prometheus.Labels{"e": "4"},
		}),
	"Io/AllowToCharge": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "io_allow_to_charge",
			Help: "",
		}),
	"Io/AllowToDischarge": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "io_allow_to_discharge",
			Help: "",
		}),
	"Io/ExternalRelay": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "io_external_relay",
			Help: "",
		}),
	"History/MinimumCellVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_min_cell_voltage_volts",
			Help: "",
		}),
	"History/MaximumCellVoltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "history_max_cell_voltage_volts",
			Help: "",
		}),
	"Pv/V": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "pv_array_voltage_volts",
			Help: "PV array voltage",
		}),
	"Pv/I": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "pv_array_current_amps",
			Help: "PV current (= /Yield/Power divided by /Pv/V)",
		}),
	"Yield/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "yield_power_watts",
			Help: "Actual input power (Watts)",
		}),
	"Yield/User": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "yield_user_total_kwh",
			Help: "Total kWh produced (user resettable)",
		}),
	"Yield/System": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "yield_system_total_kwh",
			Help: "Total kWh produced (not resettable)",
		}),
	"Load/State": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "load_state",
			Help: "Whether the load is on or off",
		}),
	"Load/I": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "load_current_amps",
			Help: "Current from the load output",
		}),
	"MppOperationMode": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "mpp_operation_mode",
			Help: "0 = Off 1 = Voltage or Current limited 2 = MPPT Tracker active",
		}),
	"Ac/Energy/Forward": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_energy_forward_kwh",
			Help: "kWh  - Total produced energy over all phases",
		}),
	"Ac/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_power_watts",
			Help: "W    - Total power of all phases, preferably real power",
		}),
	"Ac/L1/Current": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_current",
			Help:        "A AC",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/L1/Energy/Forward": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_energy_forward_kwh",
			Help:        "kWh",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/L1/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/L1/Voltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_voltage_volts",
			Help:        "V AC",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/L2/Current": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_current_amps",
			Help:        "A AC",
			ConstLabels: prometheus.Labels{"phase": "2"},
		}),
	"Ac/L2/Energy/Forward": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_energy_forward_kwh",
			Help:        "kWh",
			ConstLabels: prometheus.Labels{"phase": "2"},
		}),
	"Ac/L2/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "2"},
		}),
	"Ac/L2/Voltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_voltage_volts",
			Help:        "V AC",
			ConstLabels: prometheus.Labels{"phase": "2"},
		}),
	"Ac/L3/Current": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_current_amps",
			Help:        "A AC",
			ConstLabels: prometheus.Labels{"phase": "3"},
		}),
	"Ac/L3/Energy/Forward": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_energy_forward_kwh",
			Help:        "kWh",
			ConstLabels: prometheus.Labels{"phase": "3"},
		}),
	"Ac/L3/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "3"},
		}),
	"Ac/L3/Voltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_phase_voltage_volts",
			Help:        "V AC",
			ConstLabels: prometheus.Labels{"phase": "3"},
		}),
	"Ac/Current": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_current_amps",
			Help: "A AC - Deprecated",
		}),
	"Ac/Voltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_voltage_volts",
			Help: "V AC - Deprecated",
		}),
	"Ac/MaxPower": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_max_power_watts",
			Help: "Max rated power (in Watts) of the inverter",
		}),
	"Ac/PowerLimit": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_power_limit_watts",
			Help: "Used by the Fronius Zero-feedin feature, see ESS manual.",
		}),
	"FroniusDeviceType": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "fronius_device_type",
			Help: "Fronius specific product id list",
		}),
	"Position": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "position",
			Help: "0=AC input 1; 1=AC output; 2=AC input 2",
		}),
	"StatusCode": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "status_code",
			Help: "0=Startup 0; 1=Startup 1; 2=Startup 2; 3=Startup 4=Startup 4; 5=Startup 5; 6=Startup 6; 7=Running; 8=Standby; 9=Boot loading; 10=Error",
		}),
	"Ac/In/L1/I": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_input_phase_current_amps",
			Help:        "A AC",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/In/L1/P": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_input_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/In/CurrentLimit": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_input_current_limit_watts",
			Help: "A AC",
		}),
	"NrOfOutputs": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "output_count",
			Help: "The actual number of outputs.",
		}),
	"Dc/1/Current": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_current_amps",
			Help:        "A DC",
			ConstLabels: prometheus.Labels{"n": "1"},
		}),
	"Dc/1/Temperature": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_temperature_celcius",
			Help:        "°C - Battery temperature",
			ConstLabels: prometheus.Labels{"n": "1"},
		}),
	"Dc/2/Voltage": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_voltage_volts",
			Help:        "V DC",
			ConstLabels: prometheus.Labels{"n": "2"},
		}),
	"Dc/2/Current": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_current_amps",
			Help:        "A DC",
			ConstLabels: prometheus.Labels{"n": "2"},
		}),
	"Dc/2/Temperature": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "dc_temperature_celcius",
			Help:        "°C - Battery temperature",
			ConstLabels: prometheus.Labels{"n": "2"},
		}),
	"Ac/Energy/Reverse": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_energy_reverse_kwh",
			Help: "",
		}),
	"Ac/Grid/L1/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_grid_phase_power_watt",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/Grid/L2/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_grid_phase_power_watt",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "2"},
		}),
	"Ac/Grid/L3/Power": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_grid_phase_power_watt",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "3"},
		}),
	"Ac/Grid/NumberOfPhases": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "ac_grid_number_of_phases",
			Help: "",
		}),
	"Ac/L1/Energy/Reverse": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_energy_phase_reverse_kwh",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "1"},
		}),
	"Ac/L2/Energy/Reverse": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_energy_phase_reverse_kwh",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "2"},
		}),
	"Ac/L3/Energy/Reverse": gaugeObserver(
		prometheus.GaugeOpts{
			Name:        "ac_energy_phase_reverse_kwh",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "3"},
		}),
	"Dc/Battery/Temperature": gaugeObserver(
		prometheus.GaugeOpts{
			Name: "dc_battery_temperature_celcius",
			Help: "",
		}),
}
