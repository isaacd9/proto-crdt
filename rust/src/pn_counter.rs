use crate::{g_counter::GCounterExt, pb};
use std::{collections::HashMap, iter};

pub trait PNCounterExt {
    type T;

    fn decrement(&mut self, n: u64);
}

impl GCounterExt for pb::PnCounter {
    type T = pb::PnCounter;

    fn new(id: &'static str) -> pb::PnCounter {
        pb::PnCounter {
            identifier: id.to_string(),
            increments: HashMap::new(),
            decrements: HashMap::new(),
        }
    }

    fn increment(&mut self, n: u64) {
        self.increments
            .entry(self.identifier.to_string())
            .and_modify(|v| *v += n)
            .or_insert(n);
    }

    fn value(&self) -> u64 {
        let mut sum = 0;
        for (_, v) in &self.increments {
            sum += v
        }
        for (_, v) in &self.decrements {
            sum -= v
        }
        sum
    }

    fn merge(id: &'static str, a: &pb::PnCounter, b: &pb::PnCounter) -> pb::PnCounter {
        let mut c = Self::new(id);

        // TODO: DRY this up
        for counter in iter::once(a).chain(iter::once(b)) {
            for (id, count) in &counter.increments {
                c.increments
                    .entry(id.to_string())
                    .and_modify(|v| *v = std::cmp::max(*v, *count))
                    .or_insert(*count);
            }
        }

        for counter in iter::once(a).chain(iter::once(b)) {
            for (id, count) in &counter.decrements {
                c.decrements
                    .entry(id.to_string())
                    .and_modify(|v| *v = std::cmp::max(*v, *count))
                    .or_insert(*count);
            }
        }

        c
    }
}

impl PNCounterExt for pb::PnCounter {
    type T = pb::PnCounter;

    fn decrement(&mut self, n: u64) {
        self.decrements
            .entry(self.identifier.to_string())
            .and_modify(|v| *v += n)
            .or_insert(n);
    }
}
