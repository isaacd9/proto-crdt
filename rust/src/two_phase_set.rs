use crate::{pb, ProstMessageExt, TwoPhaseSetExt};
use bytes::Bytes;
use std::{collections::HashSet, hash::Hash};

impl<E: prost::Message + ProstMessageExt + Default + Eq + Hash> TwoPhaseSetExt<E>
    for pb::TwoPhaseSet
{
    type T = pb::TwoPhaseSet;

    fn new<I>(elements: I) -> Self::T
    where
        I: IntoIterator<Item = E>,
    {
        let any_elements = elements
            .into_iter()
            .map(|msg| prost_types::Any {
                type_url: E::type_url(),
                value: msg.encode_to_vec(),
            })
            .collect();

        pb::TwoPhaseSet {
            added: any_elements,
            removed: vec![],
        }
    }

    fn insert(&mut self, element: &E) {
        let encoded = prost_types::Any {
            type_url: E::type_url(),
            value: element.encode_to_vec(),
        };

        if !self.added.contains(&encoded) {
            self.added.push(encoded)
        }
    }

    fn remove(&mut self, element: &E) {
        let encoded = prost_types::Any {
            type_url: E::type_url(),
            value: element.encode_to_vec(),
        };

        if self.added.contains(&encoded) && !self.removed.contains(&encoded) {
            self.removed.push(encoded)
        }
    }

    fn contains(&self, element: &E) -> bool {
        let encoded = prost_types::Any {
            type_url: E::type_url(),
            value: element.encode_to_vec(),
        };

        self.added.contains(&encoded) && !self.removed.contains(&encoded)
    }

    fn len(&self) -> usize {
        self.added.len() - self.removed.len()
    }

    fn elements(&self) -> Result<HashSet<E>, prost::DecodeError> {
        let added_set: Result<HashSet<E>, prost::DecodeError> = self
            .added
            .iter()
            .map(|any| E::decode(Bytes::copy_from_slice(&any.value)))
            .collect();

        let removed_set: Result<HashSet<E>, prost::DecodeError> = self
            .removed
            .iter()
            .map(|any| E::decode(Bytes::copy_from_slice(&any.value)))
            .collect();

        let mut s = added_set?;
        for el in removed_set? {
            s.remove(&el);
        }

        Ok(s)
    }

    fn merge(a: Self::T, b: Self::T) -> Result<Self::T, prost::DecodeError> {
        let mut c = Self::T::default();

        // Added
        for a_el in a
            .added
            .into_iter()
            .map(|any| E::decode(Bytes::copy_from_slice(&any.value)))
        {
            c.insert(&a_el?);
        }
        for b_el in b
            .added
            .into_iter()
            .map(|any| E::decode(Bytes::copy_from_slice(&any.value)))
        {
            c.insert(&b_el?);
        }

        // Removed
        for a_el in a
            .removed
            .into_iter()
            .map(|any| E::decode(Bytes::copy_from_slice(&any.value)))
        {
            c.remove(&a_el?);
        }
        for b_el in b
            .removed
            .into_iter()
            .map(|any| E::decode(Bytes::copy_from_slice(&any.value)))
        {
            c.remove(&b_el?);
        }

        Ok(c)
    }
}

#[cfg(test)]
mod tests {
    #[derive(Hash, Clone, PartialEq, Eq, ::prost::Message)]
    pub struct MyProto {
        /// Identifier is a unique identifier for this replica
        #[prost(string, tag = "1")]
        pub value: ::prost::alloc::string::String,
    }

    impl crate::ProstMessageExt for MyProto {
        fn type_url() -> String {
            "type".to_string()
        }
    }

    #[test]
    fn test_two_phase_set() {
        use super::*;
        use pb::TwoPhaseSet;

        let mut a = TwoPhaseSet::new::<Vec<MyProto>>(vec![]);

        // Idempotent inserts
        a.insert(&MyProto {
            value: "hello world".to_string(),
        });
        a.insert(&MyProto {
            value: "hello world".to_string(),
        });

        // Len
        assert_eq!(1, <TwoPhaseSet as TwoPhaseSetExt<MyProto>>::len(&a));

        // Contains
        assert!(a.contains(&MyProto {
            value: "hello world".to_string()
        }));
        assert!(!a.contains(&MyProto {
            value: "bang".to_string()
        }));

        // Insert again
        a.insert(&MyProto {
            value: "bang".to_string(),
        });
        assert_eq!(2, <TwoPhaseSet as TwoPhaseSetExt<MyProto>>::len(&a));
        assert!(a.contains(&MyProto {
            value: "bang".to_string()
        }));

        // Remove
        a.remove(&MyProto {
            value: "hello world".to_string(),
        });
        assert!(!a.contains(&MyProto {
            value: "hello world".to_string()
        }));
        assert_eq!(1, <TwoPhaseSet as TwoPhaseSetExt<MyProto>>::len(&a));

        // Elements
        let set: HashSet<MyProto> = a.elements().unwrap();
        assert!(!set.contains(&MyProto {
            value: "hello world".to_string(),
        }));
        assert!(set.contains(&MyProto {
            value: "bang".to_string(),
        }));
    }
}
