
enum StatusMode {
    Pending,
    Success,
    Failed,
}

struct Task {
    TaskId: String,
    BeaconId: String,
    Command: String,
    IssuedAt: String,
    FinishedAt: String,
    Status: StatusMode,
    Output: String,
}

static mut Tasks: Vec<Task>;

fn new_task() {}
fn update_task() {}
fn get_tasks() {}
