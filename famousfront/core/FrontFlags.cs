using System.Runtime.Serialization;
using System.Windows;

namespace famousfront.core
{
  [DataContract]
  class FrontFlags
  {
    public FrontFlags()
    {
      Backend = "127.0.0.1:8002";
      KaPeriod = 100;  // milliseconds
      ContentMarginHeight = 16.0;
      FeedSourceWidth = 360.0;
      ContentMarginMinWidth = 180.0;
      ContentMaxWidth = 800d;
      TinyFontSize = 12.0;
      SmallFontSize = 13.333;
      NormalFontSize = 14.667;
      LargeFontSize = 16.0;
      VeryLargeFontSize = 17.333;
      LargestTextSize = 18.667;
      DefaultFont = new System.Windows.Media.FontFamily("Microsoft YaHei UI, Segoe UI");
      SemiBoldFont = new System.Windows.Media.FontFamily("Microsoft YaHei UI, Segoe UI Semibold");
      LightFont = new System.Windows.Media.FontFamily("Microsoft YaHei UI Light, Segoe UI Light");

      Theme = Elysium.Theme.Light;
      ThemeAccent = Elysium.AccentBrushes.Blue.Color;
      ThemeContrast = System.Windows.Media.Colors.White;
      MarginTopDown = new Thickness(0, 8, 0, 8);
      MarginAllBounds = new Thickness(4);
      EntryMargin = new Thickness(16d, 12d, 16d, 12d);
      SourceMargin = new Thickness(12d);
      ImageMaxHeight = 640.0;
      ImageMaxWidth = 400.0;
      FigureMaxWidth = 320d;
      VideoElementHeight = 331.0;
      BodyMaxWidth = 800.0;
      SearchBoxTopOffset = 64.0;
      SearchBoxMargin = new Thickness(160.0, 12, 160, 12);
      FeedSourceMargin = new Thickness(8d, 0d, 8d, 0d);
      FeedSourceFindEntriesMaxWidth = 800;
      ProgressRingLargeSize = 80.0;
      ImageTipShowDelay = 150;
      ImageTipHideDelay = 800;
      ImageTipMaxWidth = 500.0;
      FlowDocumentBlockLineHeight = 23.667d;
      ImageTipMinHeight = 600d;
      ImageTipMinWidth = 500d;
      FeedUpdateInterval = 5.0;
      ShowGalleryThreshold = 4;
    }
    [DataMember(Name = "backend")]
    public string Backend
    {
      get;
      set;
    }
    [DataMember(Name = "ka_period")]
    public int KaPeriod
    {
      get;
      set;
    }
    [DataMember(Name = "contentmargin_height")]
    public double ContentMarginHeight { get; set; }

    [DataMember(Name = "feedsource_width")]
    public double FeedSourceWidth { get; set; }

    [DataMember(Name = "contentmargin_minwidth")]
    public double ContentMarginMinWidth { get; set; }
    [DataMember(Name = "content_maxwidth")]
    public double ContentMaxWidth { get; set; }

    [DataMember(Name = "tiny_fontsize")]
    public double TinyFontSize { get; set; }
    [DataMember(Name = "small_fontsize")]
    public double SmallFontSize { get; set; }
    [DataMember(Name = "normal_fontsize")]
    public double NormalFontSize { get; set; }
    [DataMember(Name = "large_fontsize")]
    public double LargeFontSize { get; set; }
    [DataMember(Name = "verylarge_fontsize")]
    public double VeryLargeFontSize { get; set; }
    [DataMember(Name = "largest_fontsize")]
    public double LargestTextSize { get; set; }

    [DataMember(Name = "default_font")]
    public System.Windows.Media.FontFamily DefaultFont { get; set; }
    [DataMember(Name = "light_font")]
    public System.Windows.Media.FontFamily LightFont { get; set; }
    [DataMember(Name = "semibold_font")]
    public System.Windows.Media.FontFamily SemiBoldFont { get; set; }

    [DataMember(Name = "theme")]
    public Elysium.Theme Theme { get; set; }
    [DataMember(Name = "theme_accent")]
    public System.Windows.Media.Color ThemeAccent { get; set; }
    [DataMember(Name = "theme_contrast")]
    public System.Windows.Media.Color ThemeContrast { get; set; }

    [DataMember(Name = "margin_topdown")]
    public Thickness MarginTopDown { get; set; }
    [DataMember(Name = "margin_allbounds")]
    public Thickness MarginAllBounds { get; set; }
    [DataMember(Name = "image_maxheight")]
    public double ImageMaxHeight { get; set; }
    [DataMember(Name = "entry_margin")]
    public Thickness EntryMargin { get; set; }
    [DataMember(Name = "source_margin")]
    public Thickness SourceMargin { get; set; }

    [DataMember(Name = "image_maxwidth")]
    public double ImageMaxWidth { get; set; }

    [DataMember(Name = "figure_maxwidth")]
    public double FigureMaxWidth { get; set; }

    [DataMember(Name = "body_maxwidth")]
    public double BodyMaxWidth { get; set; }

    [DataMember(Name = "video_elementheight")]
    public double VideoElementHeight { get; set; }

    [DataMember(Name = "feedsource_margin")]
    public Thickness FeedSourceMargin { get; set; }
    [DataMember(Name = "searchbox_margin")]
    public Thickness SearchBoxMargin { get; set; }
    [DataMember(Name = "searchbox_topoffset")]
    public double SearchBoxTopOffset { get; set; }

    [DataMember(Name = "feedsource_findentries_maxwidth")]
    public double FeedSourceFindEntriesMaxWidth { get; set; }

    [DataMember(Name = "progressring_largesize")]
    public double ProgressRingLargeSize { get; set; }
    [DataMember(Name = "imagetip_showdelay")]
    public int ImageTipShowDelay { get; set; }// milliseconds
    [DataMember(Name = "imagetip_hidedelay")]
    public int ImageTipHideDelay { get; set; }// milliseconds
    [DataMember(Name = "imagetip_maxwidth")]
    public double ImageTipMaxWidth { get; set; }// milliseconds
    [DataMember(Name = "flowdocument_block_lineheight")]
    public double FlowDocumentBlockLineHeight { get; set; }// milliseconds
    [DataMember(Name = "imagetip_minwidth")]
    public double ImageTipMinWidth { get; set; }// milliseconds
    [DataMember(Name = "imagetip_minheight")]
    public double ImageTipMinHeight { get; set; }// milliseconds
    [DataMember(Name = "feedupdate_interval")]
    public double FeedUpdateInterval { get; set; }// minutes
    [DataMember(Name = "showgallery_threshold")]
    public int ShowGalleryThreshold { get; set; }// minutes
  }

}
