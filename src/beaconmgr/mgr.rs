use std::{error::Error, sync::Mutex};


struct Beacon {
    id: String,
    hostname: String,
    username: String,
    domain: String,
    os: String,
    arch: String,
    pid: u32,
    internal_ip: String,
    external_ip: String,
    first_seen: String,
    last_seen: String,
}

static BEACON_LIST: Mutex<Vec<Beacon>> = Mutex::new(Vec::new());
pub static CHOSEN_BEACON: Mutex<String> = Mutex::new(String::new());

pub fn new() -> Result<Beacon, Box<dyn Error>> {
    todo!()
}
pub fn mark_dead() {} 
pub fn update_beacon() {}

pub fn update_chosen_beacon(new: &String) -> Result<(), Box<dyn Error>> {
    let list = BEACON_LIST.lock().unwrap(); 
    for beacon in list.iter() {
        if beacon.id == *new {
            let mut val = CHOSEN_BEACON.lock().unwrap();
            *val = new.to_string();
            return Ok(());
        }
    }
    Err(Box::new(std::io::Error::new(
                std::io::ErrorKind::NotFound,
                "Beacon ID not found",
    )))
}
