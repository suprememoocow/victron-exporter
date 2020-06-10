package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type mqttObserver func(componentType string, componentId string, value float64)

var labels = []string{"component_type", "component_id"}

func gaugeObserver(gauge *prometheus.GaugeVec) mqttObserver {
	prometheus.MustRegister(gauge)
	return func(componentType string, componentId string, value float64) {
		gauge.WithLabelValues(componentType, componentId).Set(value)
	}
}

func alarm(alarmType string) mqttObserver {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "alarm",
			Help:        "0=OK; 1=Warning; 2=Alarm",
			ConstLabels: prometheus.Labels{"alarm_type": alarmType},
		},
		labels,
	)

	return gaugeObserver(gauge)
}

func phaseAlarm(phase string, alarmType string) mqttObserver {
	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "phase_alarm",
			Help:        "0=OK; 1=Warning; 2=Alarm",
			ConstLabels: prometheus.Labels{"phase": phase, "alarm_type": alarmType},
		},
		labels,
	)

	return gaugeObserver(gauge)
}

func counterObserver(counter *prometheus.CounterVec) mqttObserver {
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

var suffixTopicMap = map[string]mqttObserver{
	"Ac/ActiveIn/Source": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_activein_source",
			Help: "The active AC-In source of the multi",
		},
		labels,
	)),
	"Ac/Consumption/NumberOfPhases": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_consumption_number_of_phases",
			Help: "",
		},
		labels,
	)),
	"Ac/Consumption/L1/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_phase_power_watts",
			Help:        "Total of ConsumptionOnInput & ConsumptionOnOutput",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/Consumption/L2/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_phase_power_watts",
			Help:        "Total of ConsumptionOnInput & ConsumptionOnOutput",
			ConstLabels: prometheus.Labels{"phase": "2"},
		},
		labels,
	)),
	"Ac/Consumption/L3/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_phase_power_watts",
			Help:        "Total of ConsumptionOnInput & ConsumptionOnOutput",
			ConstLabels: prometheus.Labels{"phase": "3"},
		},
		labels,
	)),
	"Ac/ConsumptionOnInput/NumberOfPhases": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_consumption_on_input_number_of_phases",
			Help: "",
		},
		labels,
	)),
	"Ac/ConsumptionOnInput/L1/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_input_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/ConsumptionOnInput/L2/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_input_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "2"},
		},
		labels,
	)),
	"Ac/ConsumptionOnInput/L3/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_input_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "3"},
		},
		labels,
	)),
	"Ac/ConsumptionOnOutput/NumberOfPhases": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_consumption_on_output_number_of_phases",
			Help: "",
		},
		labels,
	)),
	"Ac/ConsumptionOnOutput/L1/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_output_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/ConsumptionOnOutput/L2/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_output_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "2"},
		},
		labels,
	)),
	"Ac/ConsumptionOnOutput/L3/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_consumption_on_output_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "3"},
		},
		labels,
	)),
	"Dc/Battery/Alarms/CircuitBreakerTripped": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_battery_alarms_circuit_breaker_tripped",
			Help: "",
		},
		labels,
	)),
	"Dc/Battery/ConsumedAmphours": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_battery_consumed_amphours",
			Help: "Ah",
		},
		labels,
	)),
	"Dc/Battery/Current": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_battery_current",
			Help: "",
		},
		labels,
	)),
	"Dc/Battery/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_battery_power_watts",
			Help: "",
		},
		labels,
	)),
	"Dc/Battery/Soc": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_battery_state_of_charge",
			Help: "",
		},
		labels,
	)),
	"Dc/Battery/State": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_battery_state",
			Help: "",
		},
		labels,
	)),
	"Dc/Battery/TimeToGo": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_battery_time_to_go_seconds",
			Help: "",
		},
		labels,
	)),
	"Dc/Battery/Voltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_battery_voltage_volts",
			Help: "",
		},
		labels,
	)),
	"Dc/Charger/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_charger_power_watts",
			Help: "",
		},
		labels,
	)),
	"Dc/Pv/Current": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_pv_current_amps",
			Help: "",
		},
		labels,
	)),
	"Dc/Pv/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_pv_power_watts",
			Help: "",
		},
		labels,
	)),
	"Dc/System/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_system_power_watts",
			Help: "",
		},
		labels,
	)),
	"Dc/Vebus/Current": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_vebus_current_amps",
			Help: "",
		},
		labels,
	)),
	"Dc/Vebus/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_vebus_power_watts",
			Help: "",
		},
		labels,
	)),
	"Buzzer/State": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "buzzer_state",
			Help: "",
		},
		labels,
	)),
	"Relay/0/State": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "relay_state",
			Help:        "",
			ConstLabels: prometheus.Labels{"relay": "0"},
		},
		labels,
	)),
	"Relay/1/State": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "relay_state",
			Help:        "",
			ConstLabels: prometheus.Labels{"relay": "1"},
		},
		labels,
	)),
	"SystemState/State": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_state",
			Help: "",
		},
		labels,
	)),
	"Timers/TimeOnGrid": counterObserver(prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "time_on_grid_seconds_total",
			Help: "Time spent on grid",
		},
		labels,
	)),
	"Timers/TimeOnGenerator": counterObserver(prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "time_on_generator_seconds_total",
			Help: "Time spent on generator",
		},
		labels,
	)),
	"Timers/TimeOnInverter": counterObserver(prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "time_on_inverter_seconds_total",
			Help: "Time spent on inverter",
		},
		labels,
	)),
	"Timers/TimeOff": counterObserver(prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "time_off_seconds_total",
			Help: "Time spent off",
		},
		labels,
	)),
	"Settings/CGwacs/AcPowerSetPoint": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_ac_power_set_point",
			Help: "User setting: Grid set-point",
		},
		labels,
	)),
	"Settings/CGwacs/BatteryLife/DischargedSoc": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_discharged_state_of_charge",
			Help: "Deprecated",
		},
		labels,
	)),
	"Settings/CGwacs/BatteryLife/DischargedTime": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_dischanged_time",
			Help: "Internal",
		},
		labels,
	)),
	"Settings/CGwacs/BatteryLife/Flags": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_flags",
			Help: "Internal",
		},
		labels,
	)),
	"Settings/CGwacs/BatteryLife/MinimumSocLimit": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_minimum_state_of_charge_limit",
			Help: "User setting: Minimum Discharge SOC",
		},
		labels,
	)),
	"Settings/CGwacs/BatteryLife/SocLimit": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_state_of_charge_limit",
			Help: "Output of the BatteryLife algorithm (read only)",
		},
		labels,
	)),
	"Settings/CGwacs/BatteryLife/State": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_battery_life_state",
			Help: "ESS state (read & write, see below)",
		},
		labels,
	)),
	"Settings/CGwacs/Hub4Mode": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_hub4_mode",
			Help: "ESS mode (read & write, see below)",
		},
		labels,
	)),
	"Settings/CGwacs/MaxChargePercentage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_max_charge_percentage",
			Help: "Deprecated",
		},
		labels,
	)),
	"Settings/CGwacs/MaxChargePower": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_max_charge_power_watts",
			Help: "User setting: Max Charge Power",
		},
		labels,
	)),
	"Settings/CGwacs/MaxDischargePercentage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_max_discharge_percentage",
			Help: "Deprecated",
		},
		labels,
	)),
	"Settings/CGwacs/MaxDischargePower": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_max_discharge_power_watts",
			Help: "User setting: Max Inverter Power",
		},
		labels,
	)),
	"Settings/CGwacs/OvervoltageFeedIn": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_overvoltage_feed_in",
			Help: "User setting: Feed-in excess solar charger power (yes/no)",
		},
		labels,
	)),
	"Settings/CGwacs/PreventFeedback": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_prevent_feedback",
			Help: "User setting: PV Inverter Zero Feed-in (on/off)",
		},
		labels,
	)),
	"Settings/CGwacs/RunWithoutGridMeter": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_cgwacs_run_without_grid_meter",
			Help: "User setting: Grid meter installed (on/off)",
		},
		labels,
	)),
	/** com.victronenergy.vebus */
	"Ac/ActiveIn/L1/F": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_active_input_phase__freq_hz",
			Help:        "Frequency",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/ActiveIn/L1/I": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_active_input_phase_current_amps",
			Help:        "Current",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/ActiveIn/L1/P": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_active_input_phase_power_watts",
			Help:        "Real power",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/ActiveIn/L1/V": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_active_input_phase_voltage_volts",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/ActiveIn/P": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_active_input_power_watts",
			Help: "Total power",
		},
		labels,
	)),
	"Ac/ActiveIn/Connected": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_active_input_connected",
			Help: "0 when inverting, 1 when connected to an AC in.",
		},
		labels,
	)),
	"Ac/ActiveIn/ActiveInput": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_active_input_active_input",
			Help: "Active input: 0 = ACin-1, 1 = ACin-2, 240 is none (inverting).",
		},
		labels,
	)),
	"Ac/In/1/CurrentLimit": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_input_current_limit",
			Help:        "",
			ConstLabels: prometheus.Labels{"input": "1"},
		},
		labels,
	)),
	"Ac/In/1/CurrentLimitIsAdjustable": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_input_current_limit_is_adjustable",
			Help:        "",
			ConstLabels: prometheus.Labels{"input": "1"},
		},
		labels,
	)),
	"Ac/In/2/CurrentLimit": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_input_current_limit",
			Help:        "",
			ConstLabels: prometheus.Labels{"input": "2"},
		},
		labels,
	)),
	"Ac/In/2/CurrentLimitIsAdjustable": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_input_current_limit_is_adjustable",
			Help:        "",
			ConstLabels: prometheus.Labels{"input": "2"},
		},
		labels,
	)),
	"Ac/PowerMeasurementType": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_power_measurement_type",
			Help: "Indicates the type of power measurement used by the system.",
		},
		labels,
	)),
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
	"Dc/0/Voltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_voltage_volts",
			Help:        "V DC",
			ConstLabels: prometheus.Labels{"n": "0"},
		},
		labels,
	)),
	"Dc/0/Current": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_current_amps",
			Help:        "A DC",
			ConstLabels: prometheus.Labels{"n": "0"},
		},
		labels,
	)),
	"Dc/0/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_power_watts",
			Help:        "",
			ConstLabels: prometheus.Labels{"n": "0"},
		},
		labels,
	)),
	"Dc/0/Temperature": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_temperature_celcius",
			Help:        "°C - Battery temperature",
			ConstLabels: prometheus.Labels{"n": "0"},
		},
		labels,
	)),
	"Mode": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mode",
			Help: "Position of the switch. 1=Charger Only;2=Inverter Only;3=On;4=Off",
		},
		labels,
	)),
	"ModeIsAdjustable": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mode_is_adjustable",
			Help: "",
		},
		labels,
	)),
	"VebusChargeState": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vebus_charge_state",
			Help: "1. Bulk, 2. Absorption, 3. Float, 4. Storage, 5. Repeat absorption, 6. Forced absorption, 7. Equalise, 8. Bulk stopped",
		},
		labels,
	)),
	"VebusSetChargeState": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "vebus_set_charge_state",
			Help: "1. Force to Equalise. 2. Force to Absorption, for maximum absorption time. 3. Force to Float, for 24 hours.",
		},
		labels,
	)),
	"Leds/Mains": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "led_mains",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		},
		labels,
	)),
	"Leds/Bulk": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "led_bulk",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		},
		labels,
	)),
	"Leds/Absorption": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "led_absoption",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		},
		labels,
	)),
	"Leds/Float": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "led_float",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		},
		labels,
	)),
	"Leds/Inverter": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "led_inverter",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		},
		labels,
	)),
	"Leds/Overload": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "led_overload",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		},
		labels,
	)),
	"Leds/LowBattery": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "led_low_battery",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		},
		labels,
	)),
	"Leds/Temperature": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "led_temperature",
			Help: "0 = Off, 1 = On, 2 = Blinking, 3 = Blinking inverted",
		},
		labels,
	)),
	/* com.victronenergy.inverter */
	"Alarms/LowVoltage":       alarm("LowVoltage"),
	"Alarms/HighVoltage":      alarm("HighVoltage"),
	"Alarms/LowTemperature":   alarm("LowTemperature"),
	"Alarms/HighTemperature":  alarm("HighTemperature"),
	"Alarms/Overload":         alarm("Overload"),
	"Alarms/LowVoltageAcOut":  alarm("LowVoltageAcOut"),
	"Alarms/HighVoltageAcOut": alarm("HighVoltageAcOut"),
	"Ac/Out/P": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_output_power_watts",
			Help: "AC Output power watts",
		},
		labels,
	)),
	"Ac/Out/L1/V": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_output_phase_volts",
			Help:        "AC Output voltage",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/Out/L1/I": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_output_phase_current_amps",
			Help:        "AC Output current",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/Out/L1/F": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_output_phase_freq_hz",
			Help:        "AC Output frequency Hertz",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/Out/L1/P": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_output_phase_power_watts",
			Help:        "Not used on vedirect inverters ",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"State": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "state",
			Help: "",
		},
		labels,
	)),
	"Dc/0/MidVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_midvoltage_volts",
			Help:        "V DC Mid voltage (BMV-702 configured to read midpoint voltage only)",
			ConstLabels: prometheus.Labels{"n": "0"},
		},
		labels,
	)),
	"Dc/0/MidVoltageDeviation": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_midvoltage_deviation_percent",
			Help:        "Percentage deviation",
			ConstLabels: prometheus.Labels{"n": "0"},
		},
		labels,
	)),
	"Dc/1/Voltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_voltage_volts",
			Help:        "V DC",
			ConstLabels: prometheus.Labels{"n": "1"},
		},
		labels,
	)),
	"ConsumedAmphours": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "consumed_amphours",
			Help: "Ah",
		},
		labels,
	)),
	"Soc": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "state_of_chare",
			Help: "0 to 100 % (BMV, BYD, Lynx BMS)",
		},
		labels,
	)),
	"TimeToGo": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "time_to_go_seconds",
			Help: "Time to in seconds (BMV SOC relay/discharge floor value, Lynx BMS).  Max value 864,000 when battery is not discharging.",
		},
		labels,
	)),
	"Info/MaxChargeCurrent": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "max_charge_current_amps",
			Help: "Charge Current Limit aka CCL  (BYD, Lynx BMS and FreedomWon)",
		},
		labels,
	)),
	"Info/MaxDischargeCurrent": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "max_discharge_current_amps",
			Help: "Discharge Current Limit aka DCL (BYD, Lynx BMS and FreedomWon)",
		},
		labels,
	)),
	"Info/MaxChargeVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "max_charge_voltage_volts",
			Help: "Maximum voltage to charge to (BYD, Lynx BMS and FreedomWon)",
		},
		labels,
	)),
	"Info/BatteryLowVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "battery_low_voltage",
			Help: "Note that Low Voltage is ignored by the system (BYD, Lynx BMS and FreedomWon)",
		},
		labels,
	)),
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
	"Settings/HasTemperature": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_has_temperature",
			Help: "",
		},
		labels,
	)),
	"Settings/HasStarterVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_has_starter_voltage",
			Help: "",
		},
		labels,
	)),
	"Settings/HasMidVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "settings_has_mid_voltage",
			Help: "",
		},
		labels,
	)),
	"History/DeepestDischarge": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_deepest_discharge",
			Help: "",
		},
		labels,
	)),
	"History/LastDischarge": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_last_discharge",
			Help: "",
		},
		labels,
	)),
	"History/AverageDischarge": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_avg_discharge",
			Help: "",
		},
		labels,
	)),
	"History/ChargeCycles": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_charge_cycles",
			Help: "",
		},
		labels,
	)),
	"History/FullDischarges": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_full_discharges",
			Help: "",
		},
		labels,
	)),
	"History/TotalAhDrawn": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_total_drawn_amphours",
			Help: "",
		},
		labels,
	)),
	"History/MinimumVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_min_voltage_volts",
			Help: "",
		},
		labels,
	)),
	"History/MaximumVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_max_voltage_volts",
			Help: "",
		},
		labels,
	)),
	"History/TimeSinceLastFullCharge": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_time_since_full_charge_seconds",
			Help: "",
		},
		labels,
	)),
	"History/AutomaticSyncs": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_automatic_syncs",
			Help: "",
		},
		labels,
	)),
	"History/LowVoltageAlarms": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_low_voltage_alarms",
			Help: "",
		},
		labels,
	)),
	"History/HighVoltageAlarms": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_high_voltage_alarms",
			Help: "",
		},
		labels,
	)),
	"History/LowStarterVoltageAlarms": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_low_starter_voltage_alarms",
			Help: "",
		},
		labels,
	)),
	"History/HighStarterVoltageAlarms": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_high_starter_voltage_alarms",
			Help: "",
		},
		labels,
	)),
	"History/MinimumStarterVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_min_starter_voltage",
			Help: "",
		},
		labels,
	)),
	"History/MaximumStarterVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_max_starter_voltage",
			Help: "",
		},
		labels,
	)),
	"History/DischargedEnergy": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_discharge_energy_kwh",
			Help: "",
		},
		labels,
	)),
	"History/ChargedEnergy": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_charged_energy_kwh",
			Help: "",
		},
		labels,
	)),
	"ErrorCode": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "error_code",
			Help: "",
		},
		labels,
	)),
	"SystemSwitch": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_switch",
			Help: "",
		},
		labels,
	)),
	"Balancing": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "balancing",
			Help: "",
		},
		labels,
	)),
	"System/NrOfBatteries": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_battery_count",
			Help: "",
		},
		labels,
	)),
	"System/BatteriesParallel": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_batteries_parallel_count",
			Help: "",
		},
		labels,
	)),
	"System/BatteriesSeries": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_batteries_series_count",
			Help: "",
		},
		labels,
	)),
	"System/NrOfCellsPerBattery": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_cells_per_battery_count",
			Help: "",
		},
		labels,
	)),
	"System/MinCellVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_min_cell_voltage_volts",
			Help: "",
		},
		labels,
	)),
	"System/MaxCellVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "system_max_cell_voltage_volts",
			Help: "",
		},
		labels,
	)),
	"Diagnostics/ShutDownsDueError": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "diagnostics_shutdowns_due_to_error_count",
			Help: "",
		},
		labels,
	)),
	"Diagnostics/LastErrors/1/Error": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "diagnostics_last_error",
			Help:        "",
			ConstLabels: prometheus.Labels{"e": "1"},
		},
		labels,
	)),
	"Diagnostics/LastErrors/2/Error": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "diagnostics_last_error",
			Help:        "",
			ConstLabels: prometheus.Labels{"e": "2"},
		},
		labels,
	)),
	"Diagnostics/LastErrors/3/Error": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "diagnostics_last_error",
			Help:        "",
			ConstLabels: prometheus.Labels{"e": "3"},
		},
		labels,
	)),
	"Diagnostics/LastErrors/4/Error": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "diagnostics_last_error",
			Help:        "",
			ConstLabels: prometheus.Labels{"e": "4"},
		},
		labels,
	)),
	"Io/AllowToCharge": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "io_allow_to_charge",
			Help: "",
		},
		labels,
	)),
	"Io/AllowToDischarge": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "io_allow_to_discharge",
			Help: "",
		},
		labels,
	)),
	"Io/ExternalRelay": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "io_external_relay",
			Help: "",
		},
		labels,
	)),
	"History/MinimumCellVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_min_cell_voltage_volts",
			Help: "",
		},
		labels,
	)),
	"History/MaximumCellVoltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "history_max_cell_voltage_volts",
			Help: "",
		},
		labels,
	)),
	"Pv/V": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pv_array_voltage_volts",
			Help: "PV array voltage",
		},
		labels,
	)),
	"Pv/I": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "pv_array_current_amps",
			Help: "PV current (= /Yield/Power divided by /Pv/V)",
		},
		labels,
	)),
	"Yield/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "yield_power_watts",
			Help: "Actual input power (Watts)",
		},
		labels,
	)),
	"Yield/User": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "yield_user_total_kwh",
			Help: "Total kWh produced (user resettable)",
		},
		labels,
	)),
	"Yield/System": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "yield_system_total_kwh",
			Help: "Total kWh produced (not resettable)",
		},
		labels,
	)),
	"Load/State": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "load_state",
			Help: "Whether the load is on or off",
		},
		labels,
	)),
	"Load/I": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "load_current_amps",
			Help: "Current from the load output",
		},
		labels,
	)),
	"MppOperationMode": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mpp_operation_mode",
			Help: "0 = Off 1 = Voltage or Current limited 2 = MPPT Tracker active",
		},
		labels,
	)),
	"Ac/Energy/Forward": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_energy_forward_kwh",
			Help: "kWh  - Total produced energy over all phases",
		},
		labels,
	)),
	"Ac/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_power_watts",
			Help: "W    - Total power of all phases, preferably real power",
		},
		labels,
	)),
	"Ac/L1/Current": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_current",
			Help:        "A AC",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/L1/Energy/Forward": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_energy_forward_kwh",
			Help:        "kWh",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/L1/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/L1/Voltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_voltage_volts",
			Help:        "V AC",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/L2/Current": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_current_amps",
			Help:        "A AC",
			ConstLabels: prometheus.Labels{"phase": "2"},
		},
		labels,
	)),
	"Ac/L2/Energy/Forward": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_energy_forward_kwh",
			Help:        "kWh",
			ConstLabels: prometheus.Labels{"phase": "2"},
		},
		labels,
	)),
	"Ac/L2/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "2"},
		},
		labels,
	)),
	"Ac/L2/Voltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_voltage_volts",
			Help:        "V AC",
			ConstLabels: prometheus.Labels{"phase": "2"},
		},
		labels,
	)),
	"Ac/L3/Current": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_current_amps",
			Help:        "A AC",
			ConstLabels: prometheus.Labels{"phase": "3"},
		},
		labels,
	)),
	"Ac/L3/Energy/Forward": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_energy_forward_kwh",
			Help:        "kWh",
			ConstLabels: prometheus.Labels{"phase": "3"},
		},
		labels,
	)),
	"Ac/L3/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "3"},
		},
		labels,
	)),
	"Ac/L3/Voltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_phase_voltage_volts",
			Help:        "V AC",
			ConstLabels: prometheus.Labels{"phase": "3"},
		},
		labels,
	)),
	"Ac/Current": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_current_amps",
			Help: "A AC - Deprecated",
		},
		labels,
	)),
	"Ac/Voltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_voltage_volts",
			Help: "V AC - Deprecated",
		},
		labels,
	)),
	"Ac/MaxPower": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_max_power_watts",
			Help: "Max rated power (in Watts) of the inverter",
		},
		labels,
	)),
	"Ac/PowerLimit": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_power_limit_watts",
			Help: "Used by the Fronius Zero-feedin feature, see ESS manual.",
		},
		labels,
	)),
	"FroniusDeviceType": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "fronius_device_type",
			Help: "Fronius specific product id list",
		},
		labels,
	)),
	"Position": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "position",
			Help: "0=AC input 1; 1=AC output; 2=AC input 2",
		},
		labels,
	)),
	"StatusCode": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "status_code",
			Help: "0=Startup 0; 1=Startup 1; 2=Startup 2; 3=Startup 4=Startup 4; 5=Startup 5; 6=Startup 6; 7=Running; 8=Standby; 9=Boot loading; 10=Error",
		},
		labels,
	)),
	"Ac/In/L1/I": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_input_phase_current_amps",
			Help:        "A AC",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/In/L1/P": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_input_phase_power_watts",
			Help:        "W",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/In/CurrentLimit": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_input_current_limit_watts",
			Help: "A AC",
		},
		labels,
	)),
	"NrOfOutputs": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "output_count",
			Help: "The actual number of outputs.",
		},
		labels,
	)),
	"Dc/1/Current": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_current_amps",
			Help:        "A DC",
			ConstLabels: prometheus.Labels{"n": "1"},
		},
		labels,
	)),
	"Dc/1/Temperature": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_temperature_celcius",
			Help:        "°C - Battery temperature",
			ConstLabels: prometheus.Labels{"n": "1"},
		},
		labels,
	)),
	"Dc/2/Voltage": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_voltage_volts",
			Help:        "V DC",
			ConstLabels: prometheus.Labels{"n": "2"},
		},
		labels,
	)),
	"Dc/2/Current": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_current_amps",
			Help:        "A DC",
			ConstLabels: prometheus.Labels{"n": "2"},
		},
		labels,
	)),
	"Dc/2/Temperature": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "dc_temperature_celcius",
			Help:        "°C - Battery temperature",
			ConstLabels: prometheus.Labels{"n": "2"},
		},
		labels,
	)),
	"Ac/Energy/Reverse": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_energy_reverse_kwh",
			Help: "",
		},
		labels,
	)),
	"Ac/Grid/L1/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_grid_phase_power_watt",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/Grid/L2/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_grid_phase_power_watt",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "2"},
		},
		labels,
	)),
	"Ac/Grid/L3/Power": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_grid_phase_power_watt",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "3"},
		},
		labels,
	)),
	"Ac/Grid/NumberOfPhases": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ac_grid_number_of_phases",
			Help: "",
		},
		labels,
	)),
	"Ac/L1/Energy/Reverse": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_energy_phase_reverse_kwh",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "1"},
		},
		labels,
	)),
	"Ac/L2/Energy/Reverse": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_energy_phase_reverse_kwh",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "2"},
		},
		labels,
	)),
	"Ac/L3/Energy/Reverse": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "ac_energy_phase_reverse_kwh",
			Help:        "",
			ConstLabels: prometheus.Labels{"phase": "3"},
		},
		labels,
	)),
	"Dc/Battery/Temperature": gaugeObserver(prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "dc_battery_temperature_celcius",
			Help: "",
		},
		labels,
	)),
}
