use crate::{pb, ProstMessageExt, TwoPhaseSet};
use bytes::Bytes;
use std::{collections::HashSet, hash::Hash};

impl<E: prost::Message + ProstMessageExt + Default + Eq + Hash> TwoPhaseSet<E> for pb::TwoPhaseSet {
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
