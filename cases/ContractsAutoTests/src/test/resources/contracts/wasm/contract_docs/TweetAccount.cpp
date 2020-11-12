#define TESTNET
// Author: zjsunzone
#include <platon/platon.hpp>
#include <string>
using namespace platon;

// data structure of a single tweet.
struct Tweet{

	public:
		uint64_t timestamp;
		std::string tweetString;

	public:
		Tweet(){}
		Tweet(uint64_t &timestamp, const std::string &tweet):timestamp(timestamp), tweetString(tweet) {
		}

	public:
		PLATON_SERIALIZE(Tweet, (timestamp)(tweetString))
};


CONTRACT TweetAccount: public platon::Contract
{
	private:
		// "array" of all tweets of this account: maps the tweet id to the actual tweet.
		platon::StorageType<"smapping"_n, std::map<std::uint64_t, Tweet>> _tweets;
		// total number of tweets in the above _tweets mapping.
		platon::StorageType<"suint"_n, uint64_t> _numberOfTweets;
		// "owner" of this account: only admin is allowed to tweet.
		platon::StorageType<"saddress"_n, Address> _adminAddress;
	
	public:
		ACTION void init()
		{
			_numberOfTweets.self() = 0;
			_adminAddress.self() = platon::platon_caller();
		}
		
		// returns true if caller of function("sender") is admin.
		CONST bool isAdmin(){
			return platon::platon_caller() == _adminAddress.self();		
		}
	
		// create new tweet
		ACTION int64_t tweet(const std::string& tweetString) {
			int64_t result = 0;		
			if(!isAdmin()){
				// only owner is allowed to create tweets for this account.
				result = -1;
			} else if (tweetString.length() > 160) {
				// tweet contains more than 160 bytes.
				result = -2;
			} else {
				_tweets.self()[_numberOfTweets].timestamp = platon_timestamp();
				_tweets.self()[_numberOfTweets].tweetString = tweetString;
				_numberOfTweets.self() = _numberOfTweets.self() + 1;	
				result = 0; // success.		
			}
			return result;
		}

		CONST std::string getTweet(uint64_t tweetId){
			// returns two values 
			std::string tweetString = _tweets.self()[tweetId].tweetString;
			uint64_t timestamp = _tweets.self()[tweetId].timestamp;
			return tweetString;		
		}
	
		CONST std::string getLatestTweet() {
			// returns three values.
			std::string tweetString = _tweets.self()[_numberOfTweets.self() - 1].tweetString;
			uint64_t timestamp = _tweets.self()[_numberOfTweets.self() - 1].timestamp;
			uint64_t numberOfTweets = _numberOfTweets.self();
			return tweetString;		
		} 
		
		CONST Address getOwnerAddress() {
			return _adminAddress.self();		
		}

		CONST uint64_t getNumberOfTweets() {
			return _numberOfTweets.self();			
		}
		
		ACTION void adminRetrieveDonations(const Address& receiver) {
			if(isAdmin()){
				Address caddr = platon_address();
				Energon e = platon_balance(caddr);
				platon_transfer(receiver, e);
			}		
		}
		
		CONST Address caddr(){
			return platon_address();		
		}
	
		CONST std::string caddrBalance(Address receiver){
			//Address caddr = platon_address();
			Energon e = platon_balance(receiver);
			return std::to_string(e.Get());		
		}
			
		ACTION void adminDeleteAccount(){
			if(isAdmin()){
				// this is a predefined function, it deletes theh contract and returns all funds to the owner.	
				platon_destroy(_adminAddress.self());	
			}		
		}


};

PLATON_DISPATCH(TweetAccount,(init)(isAdmin)(tweet)(getTweet)(getLatestTweet)
(getOwnerAddress)(getNumberOfTweets)(adminRetrieveDonations)(adminDeleteAccount)
(caddr)(caddrBalance))



