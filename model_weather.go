// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

// WeatherSection configures one weather channel, such as rain, fog, or wind.
// It defines the initial target, the channel's allowed range, and the timing and size of later transitions.
type WeatherSection struct {
	// Current sets the initial target value, the time to reach it, and how long to hold it before the channel changes again.
	Current *WeatherCurrent `json:"current,omitempty" yaml:"current,omitempty" xml:"current,omitempty"`
	// Limits sets the absolute minimum and maximum values allowed for this channel.
	Limits *WeatherMinMax `json:"limits,omitempty" yaml:"limits,omitempty" xml:"limits,omitempty"`
	// TimeLimits sets the minimum and maximum duration, in seconds, between transitions to new channel targets.
	TimeLimits *WeatherMinMax `json:"timelimits,omitempty" yaml:"timelimits,omitempty" xml:"timelimits,omitempty"`
	// ChangeLimits sets the minimum and maximum size of a change between successive channel targets.
	ChangeLimits *WeatherMinMax `json:"changelimits,omitempty" yaml:"changelimits,omitempty" xml:"changelimits,omitempty"`
	// Thresholds gates a dependent channel on the value of its controlling channel.
	// For rain and snowfall, the controlling channel is overcast.
	Thresholds *WeatherThresholds `json:"thresholds,omitempty" yaml:"thresholds,omitempty" xml:"thresholds,omitempty"`
}

// WeatherCurrent sets the initial target and transition schedule for one weather channel.
type WeatherCurrent struct {
	// Actual is the initial target value for the channel, within its Limits.
	Actual string `json:"actual,omitempty" yaml:"actual,omitempty" xml:"actual,attr,omitempty"`
	// Time is how long the channel takes to reach Actual, in seconds.
	Time string `json:"time,omitempty" yaml:"time,omitempty" xml:"time,attr,omitempty"`
	// Duration is how long the channel holds Actual after reaching it, in seconds.
	Duration string `json:"duration,omitempty" yaml:"duration,omitempty" xml:"duration,attr,omitempty"`
}

// WeatherMinMax stores lower and upper numeric bounds for a weather setting.
type WeatherMinMax struct {
	// Min is the lower bound or minimum transition value.
	Min string `json:"min,omitempty" yaml:"min,omitempty" xml:"min,attr,omitempty"`
	// Max is the upper bound or maximum transition value.
	Max string `json:"max,omitempty" yaml:"max,omitempty" xml:"max,attr,omitempty"`
}

// WeatherThresholds controls when a dependent channel may be active based on the value of its controlling channel,
// and how quickly it stops when the value leaves the allowed range.
type WeatherThresholds struct {
	// Min is the minimum controlling-channel value that permits this channel.
	Min string `json:"min,omitempty" yaml:"min,omitempty" xml:"min,attr,omitempty"`
	// Max is the maximum controlling-channel value that permits this channel.
	Max string `json:"max,omitempty" yaml:"max,omitempty" xml:"max,attr,omitempty"`
	// End is how long, in seconds, the dependent channel takes to stop after the controlling value leaves the allowed range.
	End string `json:"end,omitempty" yaml:"end,omitempty" xml:"end,attr,omitempty"`
}

// WeatherStorm configures lightning generation during storms.
type WeatherStorm struct {
	// Density controls lightning strike density on a scale from 0 to 1.
	Density string `json:"density,omitempty" yaml:"density,omitempty" xml:"density,attr,omitempty"`
	// Threshold is the overcast value, from 0 to 1, above which lightning may appear.
	Threshold string `json:"threshold,omitempty" yaml:"threshold,omitempty" xml:"threshold,attr,omitempty"`
	// Timeout is the delay between lightning strikes, in seconds.
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty" xml:"timeout,attr,omitempty"`
}
