use crate::{pb, GSetExt, ProstMessageExt};
use bytes::Bytes;
use std::collections::HashSet;
use std::hash::Hash;

impl<E: prost::Message + ProstMessageExt + Default + Eq + Hash> GSetExt<E> for pb::GSet {
    type T = pb::GSet;

    fn new<I>(elements: I) -> pb::GSet
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

    fn insert(&mut self, element: E) {
        let encoded = prost_types::Any {
            type_url: E::type_url(),
            value: element.encode_to_vec(),
        };

        if !self.elements.contains(&encoded) {
            self.elements.push(encoded)
        }
    }

    fn contains(&self, element: E) -> bool {
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

    fn merge<A, B>(a: A, b: B) -> pb::GSet
    where
        A: GSetExt<E, T = Self::T>,
        B: GSetExt<E, T = Self::T>,
    {
        pb::GSet::default()
    }
}
