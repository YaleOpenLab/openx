# CityScope II Ideas

The idea is to have multiple stakeholders and multiple tokens that interact with each other in complex ways. This notion of constant interaction requires us to rethink lots of traditional blockchain aspects because transaction fees, gas costs, etc are not necessarily needed in this model.

So we think of a new model where the blockchain we use is not strictly decentralised, has no miners because all entities need to participate in the system in order to optimize their own parameters. Since multiple people can interact at the same time and there needs to be short intervals of finality in the system, we propose a DAG based blockchain to handle latency requirements. The characteristics of the DAG based system that we want can be broadly described below:

1. Fast blocks - Block intervals should be less than 5s for quick commits of stakeholder decisions like voting on a particular scheme.
2. Immediate Finality -  Finality must be achieved within a span of 50 blocks (250s / 4 minutes) so that stakeholders can move on to vote on the next parameter.
3. Permissionless entry - anyone can enter and exit the system provided they are part of the community that is in question
4. No miners - this will be run on community hardware, so there should be no ASIC / GPU based PoW solutions
5. Faceless Leaders - Must be easy to create new decisions that people can vote on.

(These broad set of requirements are not constrained to a DAG based system, but a DAG based system works best in the situation described above)

The requirements above allow us to define the parameters of the blockchain we're using:

1. < 5s block interval time with less than 1M blocks to prevent spam
2. State commits after every x blocks, x < 300 so that we don't need to sync the blockchain from scratch every time a new node wants to join the system
3. RNG based miners - any community hw can run this
4. No block reward - if a member wants to participate in the community voting mechanism, it must download the client and sync the blockchain. It can outsource this to another party if it wants to, but that's not defined / regulated by the protocol itself
5. Easy creation of new tokens on the platform with a particular entity issuing it

While we want any entity to be able to issue a token, we could also restrict that to a particular set of users by adding in a permission based system (for eg, only the CityScope contract can add tokens and if you want to propose a new token for the cityscope contract, you need to have more than 51% of the community votes for it). This would make for autonomous governance systems as well, because the community is a pseudo DAO (pseudo since identity is tied to people within the community) and they can take decisions done by the CityScope contract that we deploy.

The characteristics of the contract itself depend on how much we want to do, to what extent, participation levels, etc but the general idea is we can regulate systems without requiring a lengthy process from stakeholders. This could also be used to do the infamous "voting on the blockchain" type of stuff, where people can choose to elect a set of community representatives and vote on them every 1/2 months regarding their efficiency. If those in power lose the election / referendum, new candidates can be proposed to take their place.

Tied in to all this is a notion / proof that a particular person is actually present in the community. The idea is that you don't need to show your proof (like SSN) directly, but only need to have a proof that you have your proof (proof of proof?). This is an interesting notion since it can be extended to many applications like decentralised identity (you only have your PoP and you can open a bank account for eg.), ownership of data and more. What matters is the way we verify the particular proof and whether we use a centralised authority there. In our application, we could use centralised parties since you are living in a piece of land that is in a sovereign state and regulated by the government. While purchasing land, there are some registrations on who lives there, who owns the property, etc, so we could use this and some other trusted data to arrive at a decentralised solution.

The extent of decentralisation and anonymity depends on the least anonymous solution, which in  this case would be the registration part. So for eg, people in a particular community would know that some person in theri community wants to have a mall near their house, needs to have 4 power sockets near their bed but wont' necessarily know who, which is the anonymous part. Since any party in the community can take part in voting, the solution is decentralised. And all states and decisions corresponding to tokens will be committed to the blockchain at frequent intervals and state history will be committed to, so there's no need to download say 100GB to take part in the system. You only need to download the last state of the system and the blocks after it and when you sync, you can delete all states that are before the most recently committed state.

Taking the raw set of ideas into consideration, we must arrive at an appropriate governance model and a suitable blockchain framework (or develop one ourselves).
