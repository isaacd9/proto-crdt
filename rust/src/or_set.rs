use crate::{pb, OrSetExt, ProstMessageExt};
use std::{collections::HashSet, hash::Hash};
use uuid::Uuid;

impl<E: prost::Message + ProstMessageExt + Default + Eq + Hash> OrSetExt<E> for pb::OrSet {
    type T = pb::OrSet;

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
            .map(|any| pb::or_set::Element {
                value: Some(any),
                identifier: Uuid::new_v4().to_string(),
            })
            .collect();

        pb::OrSet {
            added: any_elements,
            removed: vec![],
        }
    }

    fn insert(&mut self, element: &E) {
        let encoded = prost_types::Any {
            type_url: E::type_url(),
            value: element.encode_to_vec(),
        };

        let el = pb::or_set::Element {
            value: Some(encoded),
            identifier: Uuid::new_v4().to_string(),
        };
        self.added.push(el)
    }

    fn remove(&mut self, element: &E) {
        let encoded = prost_types::Any {
            type_url: E::type_url(),
            value: element.encode_to_vec(),
        };

        for el in &self.added {
            if let Some(e) = &el.value {
                if *e == encoded {
                    self.removed.push(el.clone())
                }
            }
        }
    }

    fn contains(&self, element: &E) -> bool {
        let encoded = prost_types::Any {
            type_url: E::type_url(),
            value: element.encode_to_vec(),
        };

        let mut identifiers = HashSet::new();
        for el in &self.added {
            if let Some(e) = &el.value {
                if *e == encoded {
                    identifiers.insert(&el.identifier);
                }
            }
        }
        for el in &self.removed {
            if let Some(e) = &el.value {
                if *e == encoded {
                    identifiers.remove(&el.identifier);
                }
            }
        }

        !identifiers.is_empty()
    }

    fn merge<A, B>(a: &Self::T, b: &Self::T) -> Result<Self::T, prost::DecodeError> {
        Ok(Self {
            added: a.added.iter().chain(b.added.iter()).cloned().collect(),
            removed: a.added.iter().chain(b.added.iter()).cloned().collect(),
        })
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
    fn test_or_set() {
        use super::*;
        use pb::OrSet;

        let mut a = OrSet::new::<Vec<MyProto>>(vec![]);
        // Idempotent inserts
        a.insert(&MyProto {
            value: "hello world".to_string(),
        });
        a.insert(&MyProto {
            value: "hello world".to_string(),
        });

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
    }
}
