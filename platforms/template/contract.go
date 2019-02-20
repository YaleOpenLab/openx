package template

// the contract that powers the specific platform
// the contract must at minimum have five functions similar in working to the examples below:
// 1. PreInvestmentCheck - a function which cheks if a specific investor can invest in a particular project
// 2. Invest - the main investment function which handles generation of receipt assets and similar
// 3. UpdateProjectAfterInvestment - Update the project after receipt of investor assets by the investor
// 4. sendRecipientAssets - send the recipient(s) asset(s) on acceptance of investment
// 5. UpdateProjectAfterAcceptance - update the project structure after acceptance by recipient
// Ensure all errors are caught and comments are made everywhere on how the platform functions at each stage.
// if you could draft a doc like the one in the opensolar README, it would help readers understand better
// what the platform is trying to achieve and whawt the various layers involved are.
// Also make sure you write tests so you can be sure that the platform behaves like you want it to.
