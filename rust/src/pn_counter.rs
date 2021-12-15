use crate::{pb, GCounterExt, PNCounterExt};
use std::{collections::HashMap, iter};

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
        for v in self.increments.values() {
            sum += v
        }
        for v in self.decrements.values() {
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

mod tests {
    #[test]
    fn test_pn_counter() {
        use super::*;
        use pb::PnCounter;

        let mut a = PnCounter::new("a");
        let mut b = PnCounter::new("b");
        a.increment(100);
        b.increment(200);

        a.decrement(50);
        b.decrement(30);

        let mut c = PnCounter::merge("c", &a, &b);
        assert_eq!(c.value(), 220);

        c.increment(50);
        assert_eq!(c.value(), 270);
    }
}
