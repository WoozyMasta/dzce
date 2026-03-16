// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// WeatherSection stores one weather branch configuration.
type WeatherSection struct {
	// Current defines target startup state and transition timing.
	Current *WeatherCurrent `xml:"current,omitempty" json:"current,omitempty" yaml:"current,omitempty"`
	// Limits defines absolute min/max bounds.
	Limits *WeatherMinMax `xml:"limits,omitempty" json:"limits,omitempty" yaml:"limits,omitempty"`
	// TimeLimits defines min/max random transition time window (seconds).
	TimeLimits *WeatherMinMax `xml:"timelimits,omitempty" json:"timelimits,omitempty" yaml:"timelimits,omitempty"`
	// ChangeLimits defines min/max per-step value delta window.
	ChangeLimits *WeatherMinMax `xml:"changelimits,omitempty" json:"changelimits,omitempty" yaml:"changelimits,omitempty"`
	// Thresholds defines cross-channel activation thresholds.
	// Wiki uses this mainly for rain/snowfall vs overcast dependency.
	Thresholds *WeatherThresholds `xml:"thresholds,omitempty" json:"thresholds,omitempty" yaml:"thresholds,omitempty"`
}

// WeatherCurrent stores initial weather state values.
type WeatherCurrent struct {
	// Actual is initial target value for the channel.
	Actual string `xml:"actual,attr,omitempty" json:"actual,omitempty" yaml:"actual,omitempty"`
	// Time is time to reach Actual (seconds).
	Time string `xml:"time,attr,omitempty" json:"time,omitempty" yaml:"time,omitempty"`
	// Duration is hold duration after reaching Actual (seconds).
	Duration string `xml:"duration,attr,omitempty" json:"duration,omitempty" yaml:"duration,omitempty"`
}

// WeatherMinMax stores min/max attribute pair.
type WeatherMinMax struct {
	// Min is lower bound.
	Min string `xml:"min,attr,omitempty" json:"min,omitempty" yaml:"min,omitempty"`
	// Max is upper bound.
	Max string `xml:"max,attr,omitempty" json:"max,omitempty" yaml:"max,omitempty"`
}

// WeatherThresholds stores threshold values for dependent weather channels.
type WeatherThresholds struct {
	// Min is minimum threshold.
	Min string `xml:"min,attr,omitempty" json:"min,omitempty" yaml:"min,omitempty"`
	// Max is maximum threshold.
	Max string `xml:"max,attr,omitempty" json:"max,omitempty" yaml:"max,omitempty"`
	// End is stop delay when controlling threshold is out of range (seconds).
	End string `xml:"end,attr,omitempty" json:"end,omitempty" yaml:"end,omitempty"`
}

// WeatherStorm stores storm parameters.
type WeatherStorm struct {
	// Density is lightning density (0..1).
	Density string `xml:"density,attr,omitempty" json:"density,omitempty" yaml:"density,omitempty"`
	// Threshold is overcast threshold for lightning activation (0..1).
	Threshold string `xml:"threshold,attr,omitempty" json:"threshold,omitempty" yaml:"threshold,omitempty"`
	// Timeout is delay between lightning strikes (seconds).
	Timeout string `xml:"timeout,attr,omitempty" json:"timeout,omitempty" yaml:"timeout,omitempty"`
}
