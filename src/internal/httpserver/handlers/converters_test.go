package handlers

import (
	"lo-test-task/internal/entity"
	"testing"
)

func Test_tryConvertStringToTaskStatus(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name          string
		args          args
		want          entity.TaskStatus
		isConvertable bool
	}{
		{
			name:          "\"done\" string results in entity.TaskStatusDone",
			args:          args{"done"},
			want:          entity.TaskStatusDone,
			isConvertable: true,
		},
		{
			name:          "\"in-progress\" string results in entity.TaskStatusInProgress",
			args:          args{"in-progress"},
			want:          entity.TaskStatusInProgress,
			isConvertable: true,
		},
		{
			name:          "\"created\" string results in entity.TaskStatusCreated",
			args:          args{"created"},
			want:          entity.TaskStatusCreated,
			isConvertable: true,
		},
		{
			name:          "\"incorrect\" string results in false and zero value",
			args:          args{"incorrect"},
			want:          0,
			isConvertable: false,
		},
		{
			name:          "empty string results in false and zero value",
			args:          args{""},
			want:          0,
			isConvertable: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tryConvertStringToTaskStatus(tt.args.str)
			if got != tt.want {
				t.Errorf("tryConvertStringToTaskStatus() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.isConvertable {
				t.Errorf("tryConvertStringToTaskStatus() got1 = %v, want %v", got1, tt.isConvertable)
			}
		})
	}
}

func Test_taskStatusToString(t *testing.T) {
	type args struct {
		taskStatus entity.TaskStatus
	}
	tests := []struct {
		name      string
		args      args
		want      string
		wantPanic bool
	}{
		{
			name: "entity.TaskStatusDone results in string \"done\"",
			args: args{entity.TaskStatusDone},
			want: "done",
		},
		{
			name: "entity.TaskStatusCreated results in string \"created\"",
			args: args{entity.TaskStatusCreated},
			want: "created",
		},
		{
			name: "entity.TaskStatusInProgress results in string \"in-progress\"",
			args: args{entity.TaskStatusInProgress},
			want: "in-progress",
		},
		{
			name:      "Zero value results in string panic",
			args:      args{entity.TaskStatus(0)},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if tt.wantPanic {
					if r := recover(); r == nil {
						t.Errorf("taskStatusToString() expected panic, but did not panic")
					}
				}
			}()

			if got := taskStatusToString(tt.args.taskStatus); got != tt.want {
				t.Errorf("taskStatusToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
