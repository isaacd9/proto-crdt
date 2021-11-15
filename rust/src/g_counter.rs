use std::{collections::HashMap, iter};

pub trait GCounterExt {
    type T;

    fn new(id: &'static str) -> Self::T;
    fn increment(&mut self, n: u64);
    fn value(&self) -> u64;
    fn merge(id: &'static str, a: &Self::T, b: &Self::T) -> Self::T;
}

pub mod pb {
    include!(concat!(env!("OUT_DIR"), "/pb_crdt.v1.rs"));
}

impl GCounterExt for pb::GCounter {
    type T = pb::GCounter;

    fn new(id: &'static str) -> pb::GCounter {
        pb::GCounter {
            identifier: id.to_string(),
            counts: HashMap::new(),
        }
    }

    fn increment(&mut self, n: u64) {
        self.counts
            .entry(self.identifier.to_string())
            .and_modify(|v| *v += n)
            .or_insert(n);
    }

    fn value(&self) -> u64 {
        let mut sum = 0;
        for (_, v) in &self.counts {
            sum += v
        }
        sum
    }

    fn merge(id: &'static str, a: &pb::GCounter, b: &pb::GCounter) -> pb::GCounter {
        let mut c = Self::new(id);

        for counter in iter::once(a).chain(iter::once(b)) {
            for (id, count) in &counter.counts {
                c.counts
                    .entry(id.to_string())
                    .and_modify(|v| *v = std::cmp::max(*v, *count))
                    .or_insert(*count);
            }
        }

        c
    }
}

mod tests {
    #[test]
    fn test_g_counter() {
        use super::*;
        use pb::GCounter;

        let mut a = GCounter::new("a");
        let mut b = GCounter::new("b");
        a.increment(100);
        b.increment(200);

        let mut c = GCounter::merge("c", &a, &b);
        assert_eq!(c.value(), 300);

        c.increment(50);
        assert_eq!(c.value(), 350);
    }
}
