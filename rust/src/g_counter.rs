use crate::{pb, GCounterExt, MergeExt};
use std::{collections::HashMap, iter};

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
}

impl MergeExt for pb::GCounter {
    type T = pb::GCounter;

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
