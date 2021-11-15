pub mod pb {
    include!(concat!(env!("OUT_DIR"), "/pb_crdt.v1.rs"));
}

pub mod g_counter;
pub mod pn_counter;
