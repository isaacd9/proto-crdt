use crate::{pb, OrSetExt, ProstMessageExt};
use std::{
    collections::{HashSet},
    hash::Hash,
};
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
