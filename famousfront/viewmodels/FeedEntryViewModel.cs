using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using famousfront.datamodels;
using System.Windows.Documents;
using System.Windows.Markup;
using System.Xml;
using famousfront.core;
using System.Diagnostics;
using System.Windows.Input;
using GalaSoft.MvvmLight.Command;
namespace famousfront.viewmodels
{
  class FeedEntryViewModel : famousfront.core.ViewModelBase
  {
    static readonly System.DateTime utime = new DateTime(1970, 1, 1, 0, 0, 0, 0);
    FeedEntry _ = new FeedEntry();

    internal FeedEntryViewModel(FeedEntry v)
    {
      _ = v;
      var has_summary = (_.status & FeedStatuses.Feed_status_summary_empty) == 0;
      HasDocument = (_.status & FeedStatuses.Feed_status_text_empty) == 0;
      var inline = _.status & (FeedStatuses.Feed_status_summary_mediainline | FeedStatuses.Feed_status_content_mediainline);
      var imgone = _.status & FeedStatuses.Feed_status_image_one;
      var imgmany = _.status & FeedStatuses.Feed_status_image_many;
      var video = _.videos == null || _.videos.Length < 1;
      var audio = _.audios == null || _.audios.Length < 1;
      var media = video && audio;
      if (_.images != null && _.images.Length == 1 &&_.images[0] != null && inline == 0 && media && has_summary)
      {
        Media = new ImageElementViewModel(_.images[0]);
      }
      else if (_.images != null&& _.images.Length > 1 && inline == 0 && media)
      {
        Media = new ImageGalleryViewModel(_.images);
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
        return utime.AddSeconds(_.pubdate);
      }
    }
    public bool HasMedia
    {
      get { return Media != null; }
    }
    bool is_media_inline()
    {
      var inline = _.status & (FeedStatuses.Feed_status_summary_mediainline | FeedStatuses.Feed_status_content_mediainline);
      return inline != 0;
    }
    public bool HasImageGallery
    {
      get { return !is_media_inline() && (_.status & FeedStatuses.Feed_status_image_many) != 0; }
    }
    public bool HasVideo
    {
      get { return !is_media_inline() && (_.status & (FeedStatuses.Feed_status_media_one | FeedStatuses.Feed_status_media_many)) != 0; }
    }
    public bool HasImageOne
    {
      get { return (_.status & FeedStatuses.Feed_status_image_one) != 0; }
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
      var p = utime.AddSeconds(_.pubdate);
      return new FriendlyDateTime(p);
    }
    TaskViewModel _media;
    public TaskViewModel Media
    {
      get { return _media; }
      private set { Set(ref _media, value); }
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
    public bool CanExpand
    {
      get
      {
        //var flag = FeedStatuses.Feed_status_content_empty | FeedStatuses.Feed_status_summary_empty;
        //return (_.status & flag) == 0ul;
        return true;
      }
    }
  }
  internal class FriendlyDateTime 
  {
    DateTime _;
    internal FriendlyDateTime(DateTime t)
    {
      _ = t;
    }
    public static implicit operator DateTime(FriendlyDateTime t)
    {
      return t._;
    }
    public override string ToString()
    {
      var p = _;
      var now = DateTime.Now;
      var v = p.ToString("D");
      if (p.Year != now.Year)
        return v;
      var diff = (now - p).Days;
      var dw = (int)p.DayOfWeek - 1;
      var ndw = (int)now.DayOfWeek -1;

      if (dw < 0)
        dw = 6;      
      if (ndw < 0)
        ndw = 6;

      var firstdthisweek = now.AddDays(-(int)ndw);
      var prevweek = firstdthisweek.AddDays(-7d);
      var ns = new[] { "今天", "昨天", "前天", "大前天" };
      var cws = new[] { "周一", "周二", "周三", "周四", "周五", "周六", "周日" };
      if (diff >= 0 && diff < ns.Length)
      {
        v = ns[diff];
      }else if (p >= firstdthisweek)
      {
        v = cws[dw];
      }
      else if (p >= prevweek)
      {
        v = "上" + cws[dw];
      }
      return v;
    }
  }
}
