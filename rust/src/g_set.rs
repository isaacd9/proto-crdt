use crate::{pb, GSetExt, ProstMessageExt};
use bytes::Bytes;
use std::collections::HashSet;
use std::hash::Hash;

impl<E: prost::Message + ProstMessageExt + Default + Eq + Hash> GSetExt<E> for pb::GSet {
    type T = pb::GSet;

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

        pb::GSet {
            elements: any_elements,
        }
    }

    fn insert(&mut self, element: &E) {
        let encoded = prost_types::Any {
            type_url: E::type_url(),
            value: element.encode_to_vec(),
        };

        if !self.elements.contains(&encoded) {
            self.elements.push(encoded)
        }
    }

    fn contains(&self, element: &E) -> bool {
        let encoded = prost_types::Any {
            type_url: E::type_url(),
            value: element.encode_to_vec(),
        };

        self.elements.contains(&encoded)
    }

    fn len(&self) -> usize {
        self.elements.len()
    }

    fn elements(&self) -> Result<HashSet<E>, prost::DecodeError> {
        self.elements
            .iter()
            .map(|any| E::decode(Bytes::copy_from_slice(&any.value)))
            .collect()
    }

    fn merge<A, B>(a: &A, b: &B) -> Result<Self::T, prost::DecodeError>
    where
        A: GSetExt<E>,
        B: GSetExt<E>,
    {
        let mut c = pb::GSet::default();

        for a_el in a.elements()?.into_iter() {
            c.insert(&a_el)
        }

        for b_el in b.elements()?.into_iter() {
            c.insert(&b_el)
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
    fn test_g_set() {
        use super::*;
        use pb::GSet;

        let mut a = GSet::new::<Vec<MyProto>>(vec![]);

        a.insert(&MyProto {
            value: "hello world".to_string(),
        });

        assert_eq!(1, <GSet as GSetExt<MyProto>>::len(&a));

        assert!(a.contains(&MyProto {
            value: "hello world".to_string()
        }));
        assert!(!a.contains(&MyProto {
            value: "bang".to_string()
        }));

        a.insert(&MyProto {
            value: "bang".to_string(),
        });
        assert_eq!(2, <GSet as GSetExt<MyProto>>::len(&a));
        assert!(a.contains(&MyProto {
            value: "bang".to_string()
        }));

        let set: HashSet<MyProto> = a.elements().unwrap();
        assert!(set.contains(&MyProto {
            value: "hello world".to_string(),
        }));
        assert!(set.contains(&MyProto {
            value: "bang".to_string(),
        }));
    }

    #[test]
    fn test_merge() {
        use super::*;
        use pb::GSet;

        let a = GSet::new::<Vec<MyProto>>(vec![
            MyProto {
                value: "hello".to_string(),
            },
            MyProto {
                value: "bang".to_string(),
            },
        ]);
        let b = GSet::new::<Vec<MyProto>>(vec![
            MyProto {
                value: "hello".to_string(),
            },
            MyProto {
                value: "whimper".to_string(),
            },
        ]);

        let c = <GSet as GSetExt<MyProto>>::merge(&a, &b).unwrap();
        assert!(c.contains(&MyProto {
            value: "bang".to_string()
        }));
        assert!(c.contains(&MyProto {
            value: "hello".to_string()
        }));
        assert!(c.contains(&MyProto {
            value: "whimper".to_string()
        }));

        assert_eq!(3, <GSet as GSetExt<MyProto>>::len(&c));
    }
}
