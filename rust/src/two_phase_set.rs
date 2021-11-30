use crate::{pb, OrSetExt, ProstMessageExt};
use std::hash::Hash;

impl<E: prost::Message + ProstMessageExt + Default + Eq + Hash> OrSetExt<E> for pb::TwoPhaseSet {
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
}
