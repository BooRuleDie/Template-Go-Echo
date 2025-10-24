# To Do

- [x] test storage
- [x] decide transaction flow
    * NewTX under storage
    * transaction inside the service layer
- [x] fix docker pattern in the template
- [ ] fix auth 
- [ ] decide handler > service > repo responsibilities
    * Also mention the anti-patterns for each layer
- [ ] implement embed logic
- [ ] test complex post request body input validation like the one in discover

# Errors

About error if you want to return a custom error it's handled inside the regarding layer(service, repo, handler) otherwise just return what you have to the upper layer