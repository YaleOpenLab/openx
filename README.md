# Open Financing

[![Build Status](https://travis-ci.org/OpenFinancing/openfinancing.svg?branch=master)](https://travis-ci.org/OpenFinancing/openfinancing)
[![Codecov](https://codecov.io/gh/OpenFinancing/openfinancing/branch/master/graph/badge.svg)](https://codecov.io/gh/OpenFinancing/openfinancing)

This repo contains a WIP implementation of the Open Financing model idea in stellar. Broadly, the open financing model seeks to implement the paradigm of investing without hassles and enabling smart ownership with the help of semi trusted entities on the blockchain. The Open Financing model can be thought more generally as a platform of platforms and houses multiple platforms within it (in `platforms/`).  The goal is to have a common interface (where you complete KYC, authentication, etc) and to be able to invest in multiple assets. We use the help of the blockchain to have trustless proof of ownership and debt along with a publicly auditable source of data along with proofs. Currently there are two platforms housed within OpenFinancing:

1. Housing (Bonds / Coops) - the housing project aims to make affordable housing a reality in a way that is acceptable for all stakeholders in the system (investors, residents and the community).

2. [Opensolar](https://yale-city-openlab.earth/open-microgrids) - the opensolar project aims to use schools as community centres during natural disasters like hurricanes and also aims to make schools electricity sufficient by installing solar panels on rooftop spaces. The schools themselves need not pay upfront for the solar panel cost, but instead just need to pay their electricity bill over time and through the course of payment, get ownership of the solar panels.

## FAQs

1. Why blockchain?

This is a valid question that many would have in mind. Having a decentralized, trustless and immutable form of currency like Bitcoin is *one of the applications* of a blockchain and not *the* application of blockchain systems. What we aim to do is to use the blockchain for publicly auditable proofs of ownership and data to increase overall data transparency and ease of use (investment / ownership) which the blockchain provides.

2. Why stellar?

There is no single particular reasoning. The primary reason is that it focuses towards payment systems (fast block times and finality) and the secondary reason is that they have been building on their protocol for a while. We could definitely use other solutions providing it has the advantages of stellar at the very least.

## License

[MIT](https://github.com/OpenFinancing/openfinancing/blob/master/LICENSE)
