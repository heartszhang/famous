using famousfront.datamodels;
using famousfront.messages;
using GalaSoft.MvvmLight.Command;
using System.Windows.Input;

namespace famousfront.viewmodels
{
  class MediaElementViewModel : ImageBaseViewModel
  {
    new readonly FeedMedia _;
    internal MediaElementViewModel(FeedMedia v, FeedMedia backgroundimg)
      : base(backgroundimg)
    {
      _ = v;
    }
    public string VideoUrl
    {
      get { return _.uri; }
    }
    ICommand _videoplay_command = null;
    public ICommand VideoPlayCommand
    {
      get { return _videoplay_command ?? (_videoplay_command = new RelayCommand(ExecuteVideoPlay)); }
    }
    void ExecuteVideoPlay()
    {
      MessengerInstance.Send(new VideoPlayRequest() { source = VideoUrl });
    }
  }
}
