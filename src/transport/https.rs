
enum StatusMode {
    Running,
    Stopped,
}

struct HttpsListener {
    Id: String,
    Port: u32,
    IssuedAt: String,
    StartedAt: String,
    Status: StatusMode,
    CertFile: String,
    KeyFile: String,
    Server: todo!(),
}

fn new() {}
fn start() {}
fn stop() {}

// API endpoints

// Beacon registers itself
fn post_register() {}

// Beacon posts results of task
fn post_results() {}

// Beacon polls for new tasks
fn get_task() {}
    
