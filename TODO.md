# Fountain Query


### Todo

- [ ] group_by  
- [ ] add a table where the rows are stored on disk, however updates will be sent before putting on disk, explore ideas like lru for having tables where popular rows are hot and ready in ram while others on disk  
- [ ] live aggregate functions  
- [ ] work on byte code runner (to able to run more complex logic in queries)  
- [ ] get updates from supa-base or other update streamer  

### In Progress

- [ ] joins (has full outer join, although more joins could be added)  

### Done ✓

- [x] table column validation  
- [x] stream to client  
- [x] index/channel on table  

