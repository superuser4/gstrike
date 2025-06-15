use std::{collections::HashMap, io::{self, Write}, process::exit};
use colored::Colorize;
use crate::beaconmgr;

pub fn start_repl() {
    let mut fn_map: HashMap<&str, fn(input: &Vec<String>)> = HashMap::new();
    fn_map.insert("exit", exit_func); 
    fn_map.insert("help", help_func);
    fn_map.insert("use", use_func);

    let mut input: String = String::new();
    loop {
        print!("({}) GStrike > ", beaconmgr::mgr::CHOSEN_BEACON.lock().unwrap());
        io::stdout().flush().unwrap();

        let stdin = io::stdin();
        if let Err(e) = stdin.read_line(&mut input) {
            eprintln!("{} Error when reading input: {}", "[-]".red(), e);
        }
        parse_cmd(input.trim(), &fn_map);
        input = "".to_string();
    }
}

fn parse_cmd(input: &str, fn_map: &HashMap<&str, fn(&Vec<String>)>) {
   let input_vec: Vec<String> = 
       input
       .split_whitespace()
       .map(|x| x.to_string())
       .collect();
   if input_vec.is_empty() {
       return;
   }

   if fn_map.contains_key(&input_vec[0].as_str()) {
       fn_map[&input_vec[0].as_str()](&input_vec);
   } else {
       eprintln!("No such command, try 'help'");
   }
}

fn exit_func(_input: &Vec<String>) {
    exit(0);
}

fn help_func(_input: &Vec<String>) {
   let help =
       "
       Command\t\t\tDescription
       --------\t\t\t-----------
       help\t\t\tprints help menu
       exit\t\t\tquits GStrike cli
       use\t\t\tSelects beacon for usage
       ";
   println!("{}",help);
}

fn use_func(input: &Vec<String>) {
    if let Err(e) = beaconmgr::mgr::update_chosen_beacon(&input[0]) {
        eprintln!("{} Error: {}", "[-]".red(), e);
    }
}
