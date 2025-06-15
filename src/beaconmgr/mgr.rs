
struct Beacon {
    Id: String,
    Hostname: String,
    Username: String,
    Domain: String,
    OS: String,
    Arch: String,
    PID: u32,
    InternalIP: String,
    ExternalIP: String,
    FirstSeen: String,
    LastSeen: String,
}

static mut Beacons: Vec<Beacon>;

fn new() -> Result<Beacon, std::error::Error> {}
