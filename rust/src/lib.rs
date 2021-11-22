use std::collections::HashSet;
use std::hash::Hash;

pub mod pb {
    include!(concat!(env!("OUT_DIR"), "/pb_crdt.v1.rs"));
}

pub trait GCounterExt {
    type T;

    fn new(id: &'static str) -> Self::T;
    fn increment(&mut self, n: u64);
    fn value(&self) -> u64;
    fn merge(id: &'static str, a: &Self::T, b: &Self::T) -> Self::T;
}

pub trait PNCounterExt {
    type T;

    fn decrement(&mut self, n: u64);
}

pub trait ProstMessageExt {
    fn type_url() -> String;
}

pub trait GSetExt<E: prost::Message + ProstMessageExt + Default + Eq + Hash> {
    type T;

    fn new<I, R>(elements: I) -> Self::T
    where
        I: IntoIterator<Item = E>;

    fn insert(&mut self, element: E);
    fn contains(&self, element: E) -> bool;
    fn len(&self) -> usize;

    fn elements(&self) -> Result<HashSet<E>, prost::DecodeError>;

    fn merge<A, B>(a: A, b: B) -> Result<Self::T, prost::DecodeError>
    where
        A: GSetExt<E, T = Self::T>,
        B: GSetExt<E, T = Self::T>;
}

pub mod g_counter;
pub mod g_set;
pub mod pn_counter;
