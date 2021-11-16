pub mod pb {
    include!(concat!(env!("OUT_DIR"), "/pb_crdt.v1.rs"));
}

pub trait MergeExt {
    type T;

    fn merge(id: &'static str, a: &Self::T, b: &Self::T) -> Self::T;
}
pub trait GCounterExt {
    type T;

    fn new(id: &'static str) -> Self::T;
    fn increment(&mut self, n: u64);
    fn value(&self) -> u64;
}

pub trait PNCounterExt {
    type T;

    fn decrement(&mut self, n: u64);
}

pub mod g_counter;
pub mod pn_counter;
