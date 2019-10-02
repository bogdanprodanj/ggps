package ggps

import (
	"reflect"
	"testing"
)

func TestContainsLocation(t *testing.T) {
	type args struct {
		point   []float64
		polygon [][]float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "invalid polygon size",
			args: args{
				point:   []float64{49.841993, 24.031408},
				polygon: [][]float64{{49.842180, 24.030530}, {49.842483, 24.032254}},
			},
			want: false,
		},
		{
			name: "invalid point coordinates",
			args: args{
				point:   []float64{49.841993},
				polygon: [][]float64{{49.842180, 24.030530}, {49.842483, 24.032254}, {49.841580, 24.032680}, {49.841203, 24.031001}},
			},
			want: false,
		},
		{
			name: "invalid point coordinates of the polygon",
			args: args{
				point:   []float64{49.841993, 24.031408},
				polygon: [][]float64{{49.842180, 24.030530}, {49.842483}, {49.841580, 24.032680}, {49.841203, 24.031001}},
			},
			want: false,
		},
		{
			name: "happy path",
			args: args{
				point:   []float64{49.841993, 24.031408},
				polygon: [][]float64{{49.842180, 24.030530}, {49.842483, 24.032254}, {49.841580, 24.032680}, {49.841203, 24.031001}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsLocation(tt.args.point, tt.args.polygon); got != tt.want {
				t.Errorf("ContainsLocation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDistanceBetweenPoints(t *testing.T) {
	type args struct {
		p1 []float64
		p2 []float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "invalid point coordinates",
			args: args{
				p1: []float64{},
				p2: []float64{},
			},
			want: 0,
		},
		{
			name: "happy path",
			args: args{
				p1: []float64{49.842180, 24.030530},
				p2: []float64{49.842483, 24.032254},
			},
			want: 128.2777181074759,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DistanceBetweenPoints(tt.args.p1, tt.args.p2); got != tt.want {
				t.Errorf("DistanceBetweenPoints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMidpointCoordinates(t *testing.T) {
	type args struct {
		p1 []float64
		p2 []float64
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			name: "invalid point coordinates",
			args: args{
				p1: []float64{},
				p2: []float64{},
			},
			want: nil,
		},
		{
			name: "happy path",
			args: args{
				p1: []float64{49.842180, 24.030530},
				p2: []float64{49.842483, 24.032254},
			},
			want: []float64{49.84233150319593, 24.031391997298787},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MidpointCoordinates(tt.args.p1, tt.args.p2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MidpointCoordinates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortestDistanceFromPointToLine(t *testing.T) {
	type args struct {
		o  []float64
		ab [][]float64
	}
	tests := []struct {
		name  string
		args  args
		want  float64
		want1 []float64
	}{
		{
			name: "invalid point coordinates",
			args: args{
				o:  []float64{},
				ab: [][]float64{{49.842180, 24.030530}, {49.842483, 24.032254}},
			},
			want:  0,
			want1: nil,
		},
		{
			name: "invalid line",
			args: args{
				o:  []float64{49.841993, 24.031408},
				ab: [][]float64{{49.842483, 24.032254}},
			},
			want:  0,
			want1: nil,
		},
		{
			name: "invalid line point",
			args: args{
				o:  []float64{49.841993, 24.031408},
				ab: [][]float64{{49.842180, 24.030530}, {49.842483}},
			},
			want:  0,
			want1: nil,
		},
		{
			name: "happy path",
			args: args{
				o:  []float64{49.841993, 24.031408},
				ab: [][]float64{{49.842180, 24.030530}, {49.842483, 24.032254}},
			},
			want:  36.66624539191408,
			want1: []float64{49.842310790245556, 24.031274145786774},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ShortestDistanceFromPointToLine(tt.args.o, tt.args.ab)
			if got != tt.want {
				t.Errorf("ShortestDistanceFromPointToLine() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ShortestDistanceFromPointToLine() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
