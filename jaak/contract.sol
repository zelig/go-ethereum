/// @title jaak
/// @author zelig
import "mortal";

contract Jaak is mortal {

  uint constant JAAK_FRACTION = 0.005; // percentage JAAK takes per listen
  uint constant PLAY_PRICE = 4200; /// price per play in wei

  struct Track {
    uint playCount;
    // bytes32 id;
    address addr;
  }

  // basic data structure - trackstore
  mapping (bytes32 => Track) tracks;
  // user balances
  mapping (address => uint) balances;

  /// @notice accessor for track listing
  function getTrack(bytes32 id) const returns (Track t) {
    return tracks[id];
  }

  /// @notice accessor for balances
  function getBalance(address addr) const returns (uint balance) {
    return balances[addr];
  }

  // @notice the default send is a refill, make sure register sender's balance
  function() {
    balances[msg.sender] += msg.value;
  }

  /// @notice play triggers a payment from streamer's balance to
  ///         the artists forwarding address as well as a tiny fraction
  ///         to the owner's balance
  ///
  /// @param id track ID (swarm hash of jaak manifest)
  /// @param streamer eth address of streamer
  function play(bytes32 id, address streamer) {
    uint balance = balances[streamer];
    if (balance < PLAY_PRICE) {
      throw;
    }
    track = tracks[id];
    track.playCount++;
    tracks[id] = track;
    balances[streamer] -= PLAY_PRICE;
    balances[track.addr] += (1 - JAAK_FRACTION) * PLAY_PRICE;
    balances[owner] += JAAK_FRACTION * PLAY_PRICE;
  }

  /// @notice upload creates a new track
  ///
  /// @param id track ID (swarm hash of jaak manifest)
  /// @param artist eth address of track owner (should be taken from msg.sender)
  ///        now this is managed via the jaak proxy
  function upload(bytes32 id, address artist) {
    if (msg.sender != owner) {
      throw;
    }
    tracks[id] = Track({addr: artist});
  }

  function withdraw(uint amount) {
    if (balances[msg.sender] < amount) {
      throw;
    }
    msg.sender.send(amount);
    balances[msg.sender] -= amount;
  }


}
