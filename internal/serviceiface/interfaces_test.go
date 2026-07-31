package serviceiface

import "context"

var _ UserService = (*stubUser)(nil)
var _ ProjectService = (*stubProject)(nil)
var _ TaskService = (*stubTask)(nil)
var _ CommentService = (*stubComment)(nil)
var _ NotificationService = (*stubNotification)(nil)

type stubUser struct{}

func (stubUser) Create(ctx context.Context, user User) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}
	return user, nil
}

func (stubUser) Find(ctx context.Context, id string) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}
	return User{ID: id}, nil
}

func (stubUser) FindByEmail(ctx context.Context, email string) (User, error) {
	if err := ctx.Err(); err != nil {
		return User{}, err
	}
	return User{Email: email}, nil
}

type stubProject struct{}

func (stubProject) Create(ctx context.Context, project Project) (Project, error) {
	if err := ctx.Err(); err != nil {
		return Project{}, err
	}
	return project, nil
}

func (stubProject) Find(ctx context.Context, id string) (Project, error) {
	if err := ctx.Err(); err != nil {
		return Project{}, err
	}
	return Project{ID: id}, nil
}

func (stubProject) ListByOwner(ctx context.Context, ownerID string) ([]Project, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return []Project{{OwnerID: ownerID}}, nil
}

type stubTask struct{}

func (stubTask) Create(ctx context.Context, task Task) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	return task, nil
}

func (stubTask) Find(ctx context.Context, id string) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	return Task{ID: id}, nil
}

func (stubTask) ListByProject(ctx context.Context, projectID string) ([]Task, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return []Task{{ProjectID: projectID}}, nil
}

func (stubTask) Complete(ctx context.Context, id string) (Task, error) {
	if err := ctx.Err(); err != nil {
		return Task{}, err
	}
	return Task{ID: id, Status: "complete"}, nil
}

type stubComment struct{}

func (stubComment) Create(ctx context.Context, comment Comment) (Comment, error) {
	if err := ctx.Err(); err != nil {
		return Comment{}, err
	}
	return comment, nil
}

func (stubComment) ListByTask(ctx context.Context, taskID string) ([]Comment, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return []Comment{{TaskID: taskID}}, nil
}

func (stubComment) Delete(ctx context.Context, id string) error {
	return ctx.Err()
}

type stubNotification struct{}

func (stubNotification) Notify(ctx context.Context, notification Notification) error {
	return ctx.Err()
}

func (stubNotification) ListByUser(ctx context.Context, userID string) ([]Notification, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return []Notification{{UserID: userID}}, nil
}

func (stubNotification) MarkRead(ctx context.Context, id string) error {
	return ctx.Err()
}
