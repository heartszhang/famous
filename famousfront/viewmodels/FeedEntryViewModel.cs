using System;
using famousfront.datamodels;
using famousfront.core;
using System.Windows.Input;
using famousfront.utils;
using GalaSoft.MvvmLight.Command;
namespace famousfront.viewmodels
{
  class FeedEntryViewModel : ViewModelBase
  {
    static readonly DateTime Utime = new DateTime(1970, 1, 1, 0, 0, 0, 0);
    readonly FeedEntry _ = new FeedEntry();

    internal FeedEntryViewModel(FeedEntry v)
    {
      _ = v;
      var has_summary = (_.status & FeedStatuses.FeedStatusSummaryEmpty) == 0;
      HasDocument = (_.status & FeedStatuses.FeedStatusTextEmpty) == 0;
      var inline = _.status & (FeedStatuses.FeedStatusSummaryMediainline | FeedStatuses.FeedStatusContentMediainline);
      var imgone = _.status & FeedStatuses.FeedStatusImageOne;
      var imgmany = _.status & FeedStatuses.FeedStatusImageMany;
      var video = _.videos == null || _.videos.Length < 1;
      var audio = _.audios == null || _.audios.Length < 1;
      var media = video && audio;
      if (_.images != null && _.images.Length == 1 && _.images[0] != null && inline == 0 && media && has_summary)
      {
        Media = new ImageElementViewModel(_.images[0]);
      }
      else if (_.images != null && _.images.Length > 1 && inline == 0 && media)
      {
        Media = new ImagePanelViewModel(_.images);
      }
      if (_.videos != null)
      {
        Media = new MediaElementViewModel(_.videos[0], (imgone | imgmany) != 0 ? _.images[0] : new FeedMedia());
      }
      else if (_.audios != null)
      {
        Media = new MediaElementViewModel(_.audios[0], (imgone | imgmany) != 0 ? _.images[0] : new FeedMedia());
      }

    }
    public DateTime PubDate
    {
      get
      {
        return Utime.AddSeconds(_.pubdate);
      }
    }
    public bool HasMedia
    {
      get { return Media != null; }
    }
    bool is_media_inline()
    {
      var inline = _.status & (FeedStatuses.FeedStatusSummaryMediainline | FeedStatuses.FeedStatusContentMediainline);
      return inline != 0;
    }
    public bool HasImageGallery
    {
      get { return !is_media_inline() && (_.status & FeedStatuses.FeedStatusImageMany) != 0; }
    }
    public bool HasVideo
    {
      get { return !is_media_inline() && (_.status & (FeedStatuses.FeedStatusMediaOne | FeedStatuses.FeedStatusMediaMany)) != 0; }
    }
    public bool HasImageOne
    {
      get { return (_.status & FeedStatuses.FeedStatusImageOne) != 0; }
    }
    public string Summary
    {
      get { return _.summary; }
      protected set { var p = _.summary; _.summary = value; _.content = p; RaisePropertyChanged(); }
    }
    public string Title { get { return _.title.main; } }

    public string PubDay { get { return publish_day().ToString(); } }

    FriendlyDateTime publish_day()
    {
      var p = Utime.AddSeconds(_.pubdate);
      return new FriendlyDateTime(p);
    }
    TaskViewModel _media;
    public TaskViewModel Media
    {
      get { return _media; }
      private set { Set(ref _media, value); }
    }
    bool _is_document_expanded;
    public bool IsDocumentExpanded
    {
      get { return _is_document_expanded; }
      set { Set(ref _is_document_expanded, value); }
    }
    bool _has_document = true;
    public bool HasDocument
    {
      get { return _has_document; }
      private set { Set(ref _has_document, value); }
    }
    public string Url
    {
      get { return _.uri; }
    }
    ICommand _toggle_expandsummary;
    public ICommand ToggleExpandSummaryCommand
    {
      get { return _toggle_expandsummary ?? (_toggle_expandsummary = toggle_expandsummary()); }
    }
    ICommand toggle_expandsummary()
    {
      return new RelayCommand(ExecuteToggleExpandSummary);
    }
    bool _expanded;
    public bool IsExpanded
    {
      get { return _expanded; }
      protected set { Set(ref _expanded, value); }
    }
    private void ExecuteToggleExpandSummary()
    {
      Summary = _.content;
      IsExpanded = !IsExpanded;
    }

  }
}
